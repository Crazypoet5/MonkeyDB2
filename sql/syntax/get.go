package syntax

import (
	"errors"
)

//For KV
func getParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "get" {
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
	return &SyntaxTreeNode{
		Name:      "getkv",
		Value:     tableName,
		ValueType: NAME,
		Child: []*SyntaxTreeNode{
			k,
		},
	}, nil
}
