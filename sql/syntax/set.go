package syntax

import (
	"errors"
)

//For KV
func setKVParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "set" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near : " + string(t.Raw))
	}
	tableName := t.Raw
	k, err := valueParser(tr)
	if err != nil {
		return nil, err
	}
	v, err := valueParser(tr)
	if err != nil {
		return nil, err
	}
	return &SyntaxTreeNode{
		Name:      "setkv",
		Value:     tableName,
		ValueType: NAME,
		Child: []*SyntaxTreeNode{
			k, v,
		},
	}, nil
}
