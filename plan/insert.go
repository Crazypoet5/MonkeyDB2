package plan

import (
	"../exe"
	"../sql/syntax"
	"../table"
)

func insertPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "insert" {
		return nil, nil, errors.New("Expected insert but get " + stn.Name)
	}
	r, _, err := IdenticalPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	tableName := string(r.Rows[0][0].Raw)
	table := table.GetTableByName(tableName)
	if table == nil {
		return nil, nil, errors.New("Not found this table: " + tableName)
	}
	r, _, err = fieldsPlan(stn.Child[1])
	if err != nil {
		return nil, nil, err
	}
	columnNames := make([]string, 0)
	for i := 0; i < len(r.Rows[0]); i++ {
		columnNames = append(columnNames, string(r.Rows[0][i].Raw))
	}
	//TODO: insert datas
	rel := exe.NewRelation()
	re.SetResult(len(stn.Child[2].Child))
	return rel, re, nil
}

func fieldsPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "fields" {
		return nil, nil, errors.New("Expected fields but get " + stn.Name)
	}
	row := make([]exe.Value, 0)
	for i := 0; i < len(stm.Child); i++ {
		// ignored spot case
		if stn.Child[i].Name != "identical" {
			return nil, nil, errors.New("Expected identical but get " + stn.Name)
		}
		row = append(row, stn.Child[i].Value.([]byte))
	}
	r := exe.NewRelation()
	r.AddRow(row)
	re.SetResult(0)
	return r, re, nil
}
