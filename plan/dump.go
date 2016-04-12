package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func dumpPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "dump" {
		return nil, nil, errors.New("Expected dump but get " + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("Table not found.")
	}
	reader := t.FirstPage.NewReader()
	r := reader.DumpPage()
	re.SetResult(len(r.Rows))
	return r, re, nil
}
