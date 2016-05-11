package syntax

import (
	"errors"
)

func createKVParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "createkv" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	tableName := t.Raw
	t = tr.Read()
	if t.Kind != "types" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	ktype := t.Raw
	t = tr.Read()
	if t.Kind != "types" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	vtype := t.Raw
	return &SyntaxTreeNode{
		Name:      "createkv",
		Value:     tableName,
		ValueType: NAME,
		Child: []*SyntaxTreeNode{
			&SyntaxTreeNode{
				Name:  "ktype",
				Value: ktype,
			},
			&SyntaxTreeNode{
				Name:  "vtype",
				Value: vtype,
			},
		},
	}, nil
}
