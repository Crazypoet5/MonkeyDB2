package syntax

import (
    "errors"
)

func expressionParser(tr *TokenReader) (*SyntaxTreeNode, error) {
    fork := tr.Fork()
    t := fork.Read()
    if t.Kind == "identical" {
        tr.Next(1)
        return &SyntaxTreeNode {
            Name:   "identical",
            Value:  t.Raw,
            ValueType:  NAME,
        }, nil
    }
    if t.Kind == "unReference" {
        i := fork.Read()
        if fork.Read().Kind != "unReference" {
            return nil, errors.New("You have a SQL syntax error near:" + string(i.Raw))
        }
        tr.Next(3)
        return &SyntaxTreeNode {
            Name:   "identical",
            Value:  t.Raw,
            ValueType: NAME,
        }, nil
    }
    stn, err := valueParser(tr)
    if err != nil {
        return nil, err
    }
    value, valueType := stn.Value, stn.ValueType
    return &SyntaxTreeNode {
        Name:   "identical",
        Value:  value,
        ValueType:  valueType,
    }, nil
}

func valueParser(tr *TokenReader) (*SyntaxTreeNode, error) {
    t := tr.Read()
    if t.Kind == "intval" {
        return &SyntaxTreeNode {
            Name:   "value",
            Value:  calculateInt(t.Raw),
            ValueType:  INT,
        }, nil
    }
    if t.Kind == "floatval" {
        return &SyntaxTreeNode {
            Name:   "value",
            Value:  calculateFloat(t.Raw),
            ValueType:  FLOAT,
        }, nil
    }
    return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}

func filedParser(tr *TokenReader) (*SyntaxTreeNode, error) {
    fork := tr.Fork()
    t := fork.Read()
    if t.Kind == "identical" {
        tr.Next(1)
        pre := fork.Read()
        if pre.Kind == "structs" && string(pre.Raw) == "." {
            tr.Next(1)
            snt, err := filedParser(tr)
            if err != nil {
                return nil, err
            }
            return &SyntaxTreeNode {
                Name:   "spot",
                Value:  nil,
                Child:  []*SyntaxTreeNode{
                     &SyntaxTreeNode {
                        Name:   "identcial",
                        Value:  t.Raw,
                        ValueType:  NAME,
                    },
                    snt,
                },
            }, nil
        }
        return &SyntaxTreeNode {
            Name:   "identcial",
            Value:  t.Raw,
            ValueType:  NAME,
        }, nil
    }
    return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}