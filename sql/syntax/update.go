package syntax

import (
	"errors"
)

//Unensurable
func updateParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "update" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	tName := t.Raw
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "set" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	sets, err := setsParser(tr)
	if err != nil {
		return nil, err
	}
	where, err := whereParser(tr)
	if where == nil && !tr.Empty() {
		return nil, err
	}
	return &SyntaxTreeNode{
		Name:      "update",
		Value:     tName,
		ValueType: NAME,
		Child:     []*SyntaxTreeNode{sets, where},
	}, nil
}

//Ensuable
func setParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	fork := tr.Fork()
	t := fork.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	i := t.Raw
	t = fork.Read()
	if t.Kind != "relations" || string(t.Raw) != "=" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	value, err := valueParser(fork)
	if err != nil {
		return nil, err
	}
	tr.pos = fork.pos
	return &SyntaxTreeNode{
		Name:      "set",
		Value:     i,
		ValueType: NAME,
		Child: []*SyntaxTreeNode{
			value,
		},
	}, nil
}

//Ensuable
func setsParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	sets := make([]*SyntaxTreeNode, 0)
	fork := tr.Fork()
	for {
		set, err := setParser(fork)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
		fork2 := fork.Fork()
		t := fork2.Read()
		if t.Kind == "structs" && string(t.Raw) == "," {
			fork.pos = fork2.pos
			continue
		} else {
			tr.pos = fork.pos
			return &SyntaxTreeNode{
				Name:  "sets",
				Child: sets,
			}, nil
		}
	}
}
