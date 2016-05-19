package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func showTablePlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "showtable" {
		return nil, nil, errors.New("expected showtable but get:" + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("table not exist")
	}

	rel := exe.NewRelation()
	rel.SetColumnNames([]string{"filed", "type", "key"})
	for _, v := range t.Fields {
		var key string
		if v.Index != nil {
			key = "Y"
		} else {
			key = "N"
		}
		rel.AddRow(exe.Row{exe.NewValue(exe.STRING, []byte(v.Name)),
			exe.NewValue(exe.STRING, []byte(exe.TypeToString(v.Type))),
			exe.NewValue(exe.STRING, []byte(key))})
	}

	res.SetResult(1)
	return rel, res, nil
}

func showTablesPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "showtables" {
		return nil, nil, errors.New("expected showtables but get:" + stn.Name)
	}

	rel := exe.NewRelation()
	rel.SetColumnNames([]string{"table"})
	for _, v := range table.TableList {
		rel.AddRow(exe.Row{exe.NewValue(exe.STRING, []byte(v.Name))})
	}

	res.SetResult(1)
	return rel, res, nil
}
