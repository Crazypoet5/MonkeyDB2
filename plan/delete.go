package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func deletePlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "delete" {
		return nil, nil, errors.New("expected delete but get:" + stn.Name)
	}
	if stn.Child[0].Name != "identical" {
		return nil, nil, errors.New("expected delete but get:" + stn.Name)
	}
	tName := string(stn.Child[0].Value.([]byte))
	tab := table.GetTableByName(tName)
	if tab == nil {
		return nil, nil, errors.New("Table not exists")
	}
	reader := tab.FirstPage.NewReader()
	rel := reader.DumpTable()
	if stn.Child[1] == nil {
		res.SetResult(len(rel.Rows))
		ids := &exe.BitSet{}
		for i := 0; i < len(rel.Rows); i++ {
			ids.Set(i)
		}
		tab.Delete(ids)
		return nil, res, nil
	} else {
		where, err := wherePlan(stn.Child[1])
		if err != nil {
			return nil, nil, err
		}
		ids := where(rel)
		tab.Delete(ids)
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
