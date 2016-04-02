package syntax

import (
	"errors"
)

func typeParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind == "types" {
		return &SyntaxTreeNode{
			Name:      "type",
			Value:     t.Raw,
			ValueType: NAME,
		}, nil
	}
	return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}

func columndefineParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "identical" {
		return nil, errors.New("You have a syntax error near:" + string(t.Raw))
	}
	tpNode, err := typeParser(tr)
	if err != nil {
		return nil, err
	}
	first := &SyntaxTreeNode{
		Name:  "ColumnDefine",
		Value: nil,
		Child: []*SyntaxTreeNode{
			&SyntaxTreeNode{
				Name:      "identical",
				Value:     t.Raw,
				ValueType: NAME,
			},
			tpNode,
		}}
	fork := tr.Fork()
	t2 := fork.Read()
	if t2.Kind == "attributes" {
		tr.Next(1)
		first = &SyntaxTreeNode{
			Name:      "attributes",
			Value:     t2.Raw,
			ValueType: NAME,
			Child: []*SyntaxTreeNode{
				first,
			},
		}
	}
	fork = tr.Fork()
	t2 = fork.Read()
	if t2.Kind == "structs" && string(t2.Raw) == "," {
		tr.Next(1)
		nexts, err := columndefineParser(tr)
		if err != nil {
			return nil, err
		}
		return &SyntaxTreeNode{
			Name: "dot",
			Child: []*SyntaxTreeNode{
				first, nexts,
			},
		}, nil
	}

	return first, nil
}

func createtableParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t1 := tr.Read()
	if t1.Kind != "keyword" || string(t1.Raw) != "create" {
		return nil, errors.New("You have a syntax error near:" + string(t1.Raw))
	}
	t2 := tr.Read()
	if t2.Kind != "keyword" || string(t2.Raw) != "table" {
		return nil, errors.New("You have a syntax error near:" + string(t2.Raw))
	}
	t3 := tr.Read()
	if t3.Kind != "identical" {
		return nil, errors.New("You have a syntax error near:" + string(t3.Raw))
	}
	t4 := tr.Read()
	if t4.Kind != "structs" || string(t4.Raw) != "(" {
		return nil, errors.New("You have a syntax error near:" + string(t4.Raw))
	}
	columndefine, err := columndefineParser(tr)
	if err != nil {
		return nil, err
	}
	t5 := tr.Read()
	if t5.Kind != "structs" || string(t5.Raw) != ")" {
		return nil, errors.New("You have a syntax error near:" + string(t5.Raw))
	}
	return &SyntaxTreeNode{
		Name: "createtable",
		Child: []*SyntaxTreeNode{
			&SyntaxTreeNode{
				Name:      "identical",
				ValueType: NAME,
				Value:     t3.Raw,
			},
			columndefine,
		},
	}, nil
}
