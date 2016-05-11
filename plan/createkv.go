package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func createKVPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	result := NewResult()
	if stn.Name != "createkv" {
		return nil, nil, errors.New("expected createkv but get:" + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(string(tableName))
	if t != nil {
		return nil, nil, errors.New("table already exists")
	}
	ktypeR := stn.Child[0].Value.([]byte)
	ktype := exe.StringToType(string(ktypeR))
	vtypeR := stn.Child[1].Value.([]byte)
	vtype := exe.StringToType(string(vtypeR))
	ret := table.CreateKVTable(tableName, ktype, vtype)
	if ret == nil {
		return nil, nil, errors.New("create table failed")
	}
	result.SetResult(0)
	return nil, result, nil
}
