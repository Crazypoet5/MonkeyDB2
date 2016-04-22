package syntax

import (
	"errors"
)

func dropIndexParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "drop" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "index" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	return &SyntaxTreeNode{
		Name:      "dropindex",
		Child:     nil,
		Value:     t.Raw,
		ValueType: NAME,
	}, nil
}
