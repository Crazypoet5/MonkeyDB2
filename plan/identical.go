package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
)

func IdenticalPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	result := NewResult()
	if stn.Name != "identical" {
		return nil, nil, errors.New("Expect indentical, but get " + stn.Name)
	}
	relation := exe.NewRelation()
	value := exe.NewValue(exe.STRING, stn.Value.([]byte))
	row := exe.NewRow([]exe.Value{value})
	relation.AddRow(row)
	result.SetResult(1)
	return relation, result, nil
}
