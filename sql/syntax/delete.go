package syntax

import (
	"errors"
)

func deleteParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" && string(t.Raw) != "delete" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "from" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	tName := tr.Read()
	if tName.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(tName.Raw))
	}
	where, err := whereParser(tr)
	if where == nil && !tr.Empty() {
		return nil, err
	}
	return &SyntaxTreeNode{
		Name: "delete",
		Child: []*SyntaxTreeNode{&SyntaxTreeNode{
			Name:      "identical",
			Value:     tName.Raw,
			ValueType: NAME,
		}, where},
	}, nil
}
