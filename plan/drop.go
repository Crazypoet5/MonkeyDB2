package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func dropPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "drop" {
		return nil, nil, errors.New("Expected drop but get " + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("Table not found.")
	}
	t.Drop()
	re.SetResult(0)
	return nil, re, nil
}
