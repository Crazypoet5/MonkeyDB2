package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
)

func DirectPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	if stn.Name == "createtable" {
		return CreatePlan(stn)
	}
	if stn.Name == "insert" {
		return insertPlan(stn)
	}
	if stn.Name == "dump" {
		return dumpPlan(stn)
	}
	if stn.Name == "select" {
		return selectPlan(stn)
	}
	return nil, nil, errors.New("Unsopprted plan.")
}
