package syntax

import (
	"errors"
)

func dumpParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "dump" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	return &SyntaxTreeNode{
		Name:      "dump",
		Child:     nil,
		Value:     t.Raw,
		ValueType: NAME,
	}, nil
}
