package syntax

import (
	"errors"
)

func valueParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	fork := tr.Fork()
	t := fork.Read()
	if t.Kind == "intval" {
		tr.Next(1)
		return &SyntaxTreeNode{
			Name:      "value",
			Value:     calculateInt(t.Raw),
			ValueType: INT,
		}, nil
	}
	if t.Kind == "floatval" {
		tr.Next(1)
		return &SyntaxTreeNode{
			Name:      "value",
			Value:     calculateFloat(t.Raw),
			ValueType: FLOAT,
		}, nil
	}
	if t.Kind == "string" {
		str, err := stringParser(tr)
		if err != nil {
			return nil, err
		}
		return str, nil
	}
	return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}
