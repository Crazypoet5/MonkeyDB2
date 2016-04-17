package syntax

import (
	"errors"
)

//Unensurable
func selectParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "select" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	projects, err := selectFiledsParser(tr)
	if err != nil {
		return nil, err
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "from" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	froms, err := selectFiledsParser(tr)
	if err != nil {
		return nil, err
	}
	where, err := whereParser(tr)
	if where == nil && !tr.Empty() {
		return nil, err
	}
	return &SyntaxTreeNode{
		Name:  "select",
		Child: []*SyntaxTreeNode{projects, froms, where},
	}, nil
}

//Unensurable
func selectFiledsParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	vs := make([]*SyntaxTreeNode, 0)
	for {
		v, err := filedParser(tr)
		if err == nil {
			vs = append(vs, v)
		} else {
			return nil, err
		}
		fork := tr.Fork()
		t := fork.Read()
		if t.Kind == "structs" && string(t.Raw) == "," {
			tr.Next(1)
			continue
		} else {
			break
		}
	}
	return &SyntaxTreeNode{
		Name:  "fields",
		Child: vs,
	}, nil
}
