package plan

import (
	"errors"
	"unsafe"

	"../exe"
	"../sql/syntax"
	"../table"
)

//Unensuarble
func insertPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "insert" {
		return nil, nil, errors.New("Expected insert but get " + stn.Name)
	}
	r, _, err := IdenticalPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	tableName := string(r.Rows[0][0].Raw)
	table := table.GetTableByName(tableName)
	if table == nil {
		return nil, nil, errors.New("Not found this table: " + tableName)
	}
	r, _, err = fieldsPlan(stn.Child[1])
	if err != nil {
		return nil, nil, err
	}
	columnNames := make([]string, 0)
	for i := 0; len(r.Rows) > 0 && i < len(r.Rows[0]); i++ {
		columnNames = append(columnNames, string(r.Rows[0][i].Raw))
	}

	datas := make([][][]byte, 0)
	r, _, err = rowsPlan(stn.Child[2])
	for _, row := range r.Rows {
		rowD := make([][]byte, 0)
		for i := 0; i < len(row); i++ {
			rowD = append(rowD, row[i].Raw)
		}
		datas = append(datas, rowD)
	}
	if len(columnNames) == 0 {
		for i := 0; i < len(datas[0]); i++ {
			columnNames = append(columnNames, table.Fields[i].Name)
		}
	}
	for i := 0; i < len(datas); i++ {
		if len(columnNames) != len(datas[i]) {
			return nil, nil, errors.New("Insert fields does not match data row")
		}
	}
	fieldMap := make(map[int]int)
	n := 0
	for fk, f := range table.Fields {
		for dk, v := range columnNames {
			if v == f.Name {
				fieldMap[fk] = dk
				n++
				continue
			}
		}
	}
	for k, f := range table.Fields {
		if _, ok := fieldMap[k]; f.Index != nil && !ok {
			return nil, nil, errors.New(f.Name + " field is a index, but you give no value.")
		}
	}
	if n != len(columnNames) {
		return nil, nil, errors.New("Insert fields error, maybe repeat or not exist in table")
	}
	err = table.Insert(fieldMap, datas)
	if err != nil {
		return nil, nil, err
	}
	rel := exe.NewRelation()
	re.SetResult(len(stn.Child[2].Child))
	return rel, re, nil
}

//Unensurable
func fieldsPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	r := exe.NewRelation()
	if stn == nil {
		re.SetResult(0)
		return r, re, nil
	}
	if stn.Name != "fields" {
		return nil, nil, errors.New("Expected fields but get " + stn.Name)
	}
	row := make([]exe.Value, 0)
	for i := 0; i < len(stn.Child); i++ {
		// ignored spot case
		if stn.Child[i].Name != "identical" {
			return nil, nil, errors.New("Expected identical but get " + stn.Child[i].Name)
		}
		row = append(row, exe.NewValue(exe.STRING, stn.Child[i].Value.([]byte)))
	}
	r.AddRow(row)
	re.SetResult(0)
	return r, re, nil
}

//Unensurable
func rowsPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "rows" {
		return nil, nil, errors.New("Expected rows but get " + stn.Name)
	}
	r := exe.NewRelation()
	for rowIndex := 0; rowIndex < len(stn.Child); rowIndex++ {
		rowNode := stn.Child[rowIndex]
		if rowNode.Name != "row" {
			return nil, nil, errors.New("Expected row but get " + stn.Name)
		}
		var row exe.Row
		for i := 0; i < len(rowNode.Child); i++ {
			v := rowNode.Child[i]
			switch v.Name {
			case "value":
				var b *[8]byte
				var kind int
				if v.ValueType == syntax.INT {
					i := v.Value.(int)
					b = (*[8]byte)(unsafe.Pointer(&i))
					kind = exe.FLOAT
				} else {
					f := v.Value.(float64)
					b = (*[8]byte)(unsafe.Pointer(&f))
					kind = exe.INT
				}
				var bs []byte
				for i := 0; i < 8; i++ {
					bs = append(bs, (*b)[i])
				}
				row = append(row, exe.NewValue(kind, bs))
			case "string":
				row = append(row, exe.NewValue(exe.STRING, v.Value.([]byte)))
			}
		}
		r.AddRow(row)
	}

	re.SetResult(1)
	return r, re, nil
}

func uint322bytes(u uint32) []byte {
	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b[i] = (byte)(u)
		u >>= 8
	}
	return b
}
