package plan

import (
	"errors"

	"../sql/syntax"
)

func StringPlan(stn *syntax.SyntaxTreeNode) (*Relation, *Result, error) {
	result := NewResult()
	if stn.Name != "string" {
		return nil, nil, errors.New("Expect string, but get " + stn.Name)
	}
	relation := NewRelation()
	value := NewValue(STRING, stn.Value.([]byte))
	row := NewRow([]Value{value})
	relation.AddRow(row)
	result.SetResult(1)
	return relation, result, nil
}
