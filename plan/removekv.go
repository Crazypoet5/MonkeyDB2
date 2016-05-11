package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func removeKVPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "removekv" {
		return nil, nil, errors.New("expected removekv but get:" + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("table not exist")
	}
	k := stn.Child[0].Value.([]byte)
	err := t.KVRemove(k)
	if err != nil {
		return nil, nil, err
	}
	res.SetResult(1)
	return nil, res, nil
}
