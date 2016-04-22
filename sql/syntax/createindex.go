package syntax

import (
	"errors"
)

func createIndexParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "create" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "index" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	indexName := t.Raw
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "on" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	tableName := t.Raw
	t = tr.Read()
	if t.Kind != "structs" || string(t.Raw) != "(" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	fieldName := t.Raw
	t = tr.Read()
	if t.Kind != "structs" || string(t.Raw) != ")" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	return &SyntaxTreeNode{
		Name:  "createindex",
		Value: [][]byte{indexName, tableName, fieldName},
	}, nil
}
