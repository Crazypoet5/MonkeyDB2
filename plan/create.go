package plan

import (
	"errors"
	"strconv"

	"../exe"

	"../index"
	"../table"

	"../sql/syntax"
)

func CreatePlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	result := NewResult()
	if stn.Name != "createtable" {
		return nil, nil, errors.New("Expected createtable but get " + stn.Name)
	}
	r, _, err := IdenticalPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	tableName := string(r.Rows[0][0].Raw)
	t := table.CreateTable(tableName)
	r, _, err = ColumnDefinesPlan(stn.Child[1])
	if err != nil {
		return nil, nil, err
	}
	for _, v := range r.Rows {
		tp, _ := strconv.Atoi(string(v[1].Raw))
		size := 0
		fixed := false
		switch tp {
		case exe.INT:
			size = 8
		case exe.FLOAT:
			size = 8
		case exe.STRING:
			fixed = false
		case exe.ARRAY:
			fixed = false
		case exe.OBJECT:
			size = 8
		}
		t.AddFiled(string(v[0].Raw), fixed, size, tp, index.PRIMARY)
	}

	result.SetResult(0)
	return nil, result, nil
}

func ColumnDefinesPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	result := NewResult()
	relation := exe.NewRelation()
	switch stn.Name {
	case "ColumnDefine":
		if stn.Child[0].Name != "identical" {
			return nil, nil, errors.New("Expected indentical but get " + stn.Child[0].Name)
		}
		varName := stn.Child[0].Value.([]byte)
		if stn.Child[1].Name != "type" {
			return nil, nil, errors.New("Expected type but get " + stn.Child[1].Name)
		}
		varType := exe.StringToType(string(stn.Child[1].Value.([]byte)))
		row := exe.NewRow([]exe.Value{exe.NewValue(exe.STRING, varName), exe.NewValue(exe.INT, []byte(strconv.Itoa(varType)))})
		relation.AddRow(row)
		result.SetResult(1)
		return relation, result, nil
	case "attributes":
		relation, r, err := ColumnDefinesPlan(stn.Child[0])
		if err != nil {
			return nil, nil, err
		}
		result.SetResult(r.AffectedRows)
		return relation, result, nil
	case "dot":
		r, re, err := ColumnDefinesPlan((stn.Child[0]))
		if err != nil {
			return nil, nil, err
		}
		Num := re.AffectedRows
		Rows := r.Rows
		r, re, err = ColumnDefinesPlan((stn.Child[1]))
		if err != nil {
			return nil, nil, err
		}
		Num += re.AffectedRows
		Rows = append(Rows, r.Rows...)
		for _, v := range Rows {
			relation.AddRow(v)
		}
		result.SetResult(Num)
		return relation, result, nil
	}
	return nil, nil, errors.New("Cannot plan " + stn.Name)
}
