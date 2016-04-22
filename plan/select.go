package plan

import (
	"errors"

	"../exe"
	"../sql/syntax"
	"../table"
)

func selectPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "select" {
		return nil, nil, errors.New("expected select but get:" + stn.Name)
	}
	projects, err := selectFieldsPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	froms, err := selectFieldsPlan(stn.Child[1])
	if err != nil {
		return nil, nil, err
	}
	var tab *table.Table
	if len(froms) == 1 {
		tab = table.GetTableByName(froms[0])
	} else {
		//TODO: Join
	}
	if tab == nil {
		return nil, nil, errors.New("Table not found.")
	}
	if len(projects) == 1 && projects[0] == "*" {
		projects = []string{}
		for i := 0; i < len(tab.Fields); i++ {
			projects = append(projects, tab.Fields[i].Name)
		}
	}
	reader := tab.FirstPage.NewReader()
	rel := reader.DumpTable()
	if stn.Child[2] == nil {
		res.SetResult(len(rel.Rows))
		return rel.Project(projects), res, nil
	} else {
		where, err := wherePlan(stn.Child[2])
		if err != nil {
			return nil, nil, err
		}
		ids := where(rel)
		relN := exe.NewRelation()
		relN.SetColumnNames(rel.ColumnNames)
		for i := 0; i < ids.Len(); i++ {
			if ids.Get(i) {
				relN.AddRow(rel.Rows[i])
			}
		}
		res.SetResult(len(relN.Rows))
		return relN.Project(projects), res, nil
	}
}

func selectFieldsPlan(stn *syntax.SyntaxTreeNode) ([]string, error) {
	if stn.Name != "fields" {
		return nil, errors.New("expected fields but get:" + stn.Name)
	}
	ret := make([]string, 0)
	for i := 0; i < len(stn.Child); i++ {
		ret = append(ret, string(stn.Child[i].Value.([]byte)))
	}
	return ret, nil
}
