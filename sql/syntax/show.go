package syntax

import (
	"errors"
)

func showParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "show" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "table" {
		if t.Kind == "keyword" && string(t.Raw) == "tables" {
			return &SyntaxTreeNode{
				Name: "showtables",
			}, nil
		}
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	tableName := t.Raw
	return &SyntaxTreeNode{
		Name:      "showtable",
		Value:     tableName,
		ValueType: NAME,
	}, nil
}
