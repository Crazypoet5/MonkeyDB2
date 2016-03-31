package plan

import (
	"errors"
	"fmt"
	"strconv"

	"../sql/syntax"
)

func CreatePlan(stn *syntax.SyntaxTreeNode) (*Relation, *Result, error) {
	result := NewResult()
	if stn.Name != "createtable" {
		return nil, nil, errors.New("Expected createtable but get " + stn.Name)
	}
	r, _, err := StringPlan(stn.Child[0])
	if err != nil {
		return nil, nil, err
	}
	tableName := string(r.Rows[0][0].Raw)
	fmt.Println(tableName)
	r, _, err = ColumnDefinesPlan(stn.Child[1])
	if err != nil {
		return nil, nil, err
	}
	for _, v := range r.Rows {
		fmt.Println(v)
	}
	result.SetResult(0)
	return nil, result, nil
}

func ColumnDefinesPlan(stn *syntax.SyntaxTreeNode) (*Relation, *Result, error) {
	result := NewResult()
	relation := NewRelation()
	switch stn.Name {
	case "ColumnDefine":
		if stn.Child[0].Name != "identical" {
			return nil, nil, errors.New("Expected indentical but get " + stn.Child[0].Name)
		}
		varName := stn.Child[0].Value.([]byte)
		if stn.Child[1].Name != "type" {
			return nil, nil, errors.New("Expected type but get " + stn.Child[1].Name)
		}
		varType := StringToType(string(stn.Child[1].Value.([]byte)))
		row := NewRow([]Value{NewValue(STRING, varName), NewValue(INT, []byte(strconv.Itoa(varType)))})
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
