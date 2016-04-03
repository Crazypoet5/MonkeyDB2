package syntax

import (
	"errors"
)

func insertParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "insert" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	t = tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "into" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	tName := tr.Read()
	if tName.Kind != "identical" {
		return nil, errors.New("You have a syntax error near: " + string(tName.Raw))
	}
	fork := tr.Fork()
	t = fork.Read()
	var names *SyntaxTreeNode
	if t.Kind == "structs" && string(t.Raw) == "(" {
		var err error
		names, err = filedsParser(tr)
		if err != nil {
			return nil, err
		}
	}
	rows, err := valuesParser(tr)
	if err != nil {
		return nil, err
	}
	return &SyntaxTreeNode{
		Name: "insert",
		Child: []*SyntaxTreeNode{
			&SyntaxTreeNode{
				Name:      "identical",
				Value:     tName.Raw,
				ValueType: NAME,
			},
			names,
			rows,
		},
	}, nil
}

func valuesParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "keyword" || string(t.Raw) != "values" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	vs, err := rowsParser(tr)
	if err != nil {
		return nil, err
	}
	return vs, nil
}

func rowParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "structs" || string(t.Raw) != "(" {
		return nil, errors.New("You have a syntax error near: " + string(t.Raw))
	}
	vs := make([]*SyntaxTreeNode, 0)
	for {
		v, err := valueParser(tr)
		if err == nil {
			vs = append(vs, v)
		} else {
			return nil, err
		}
		t := tr.Read()
		if t.Kind == "structs" && string(t.Raw) == "," {
			continue
		} else if t.Kind == "structs" && string(t.Raw) == ")" {
			break
		} else {
			return nil, errors.New("You have a syntax error near " + string(t.Raw))
		}
	}
	return &SyntaxTreeNode{
		Name:  "row",
		Child: vs,
	}, nil
}

func rowsParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	vs := make([]*SyntaxTreeNode, 0)
	v, err := rowParser(tr)
	if err != nil {
		return nil, err
	}
	vs = append(vs, v)
	for {
		fork := tr.Fork()
		t := fork.Read()
		if t.Kind != "structs" || string(t.Raw) != "," {
			break
		}
		tr.Next(1)
		v, err = rowParser(tr)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return &SyntaxTreeNode{
		Name:  "rows",
		Child: vs,
	}, nil
}
