package plan

import (
	"errors"
	"unsafe"

	"../exe"
	"../sql/syntax"
	"../table"
)

func updatePlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "update" {
		return nil, nil, errors.New("Except update but get: " + stn.Name)
	}
	tabN := string(stn.Value.([]byte))
	tab := table.GetTableByName(tabN)
	if tab == nil {
		return nil, nil, errors.New("Table not exists")
	}
	sets, err := setsPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	n := 0
	for _, v := range tab.Fields {
		if _, ok := sets[v.Name]; ok {
			n++
		}
	}
	if n != len(sets) {
		return nil, nil, errors.New("Some fields given are not in table.")
	}
	reader := tab.FirstPage.NewReader()
	rel := reader.DumpTable()
	if stn.Child[1] == nil {
		res.SetResult(len(rel.Rows))
		ids := &exe.BitSet{}
		for i := 0; i < len(rel.Rows); i++ {
			ids.Set(i)
		}
		newColumns := make(map[int]int)
		for i := 0; i < len(rel.ColumnNames); i++ {
			newColumns[i] = i
		}
		var data [][][]byte
		for r := 0; r < len(rel.Rows); r++ {
			data = append(data, make([][]byte, 0))
			for i := 0; i < len(rel.ColumnNames); i++ {
				if d, ok := sets[rel.ColumnNames[i]]; ok {
					data[r] = append(data[r], d)
				} else {
					data[r] = append(data[r], rel.Rows[r][i].Raw)
				}
			}
		}

		tab.Delete(ids)
		tab.Insert(newColumns, data)
		return nil, res, nil
	} else {
		where, err := wherePlan(stn.Child[1])
		if err != nil {
			return nil, nil, err
		}
		ids := where(rel)
		newColumns := make(map[int]int)
		for i := 0; i < len(rel.ColumnNames); i++ {
			newColumns[i] = i
		}
		var data [][][]byte
		for r := 0; r < len(rel.Rows); r++ {
			if ids.Get(r) {
				data = append(data, make([][]byte, 0))
				for i := 0; i < len(rel.ColumnNames); i++ {
					if d, ok := sets[rel.ColumnNames[i]]; ok {
						data[len(data)-1] = append(data[len(data)-1], d)
					} else {
						data[len(data)-1] = append(data[len(data)-1], rel.Rows[r][i].Raw)
					}
				}
			}

		}
		tab.Delete(ids)
		tab.Insert(newColumns, data)
		n := 0
		for i := 0; i < ids.Len(); i++ {
			if ids.Get(i) {
				n++
			}
		}
		res.SetResult(n)
		return nil, res, nil
	}
}

func setsPlan(stn *syntax.SyntaxTreeNode) (map[string][]byte, error) {
	if stn.Name != "sets" {
		return nil, errors.New("Except sets but get: " + stn.Name)
	}
	ret := make(map[string][]byte)
	for _, v := range stn.Child {
		if v.Name != "set" {
			return nil, errors.New("Except set but get: " + v.Name)
		}
		col := string(v.Value.([]byte))
		switch v.Child[0].Name {
		case "value":
			var b *[8]byte
			if v.Child[0].ValueType == syntax.INT {
				i := v.Child[0].Value.(int)
				b = (*[8]byte)(unsafe.Pointer(&i))
			} else {
				f := v.Child[0].Value.(float64)
				b = (*[8]byte)(unsafe.Pointer(&f))
			}
			var bs []byte
			for i := 0; i < 8; i++ {
				bs = append(bs, (*b)[i])
			}
			ret[col] = bs
		case "string":
			ret[col] = v.Child[0].Value.([]byte)
		}
	}
	return ret, nil
}
