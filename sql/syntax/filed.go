package syntax

import (
    "errors"
)

func expressionParser(tr *TokenReader) (*syntaxTreeNode, error) {
    fork := tr.Fork()
    t := fork.Read()
    if t.Kind == "identical" {
        tr.Next(1)
        return &syntaxTreeNode {
            name:   "identical",
            value:  t.Raw,
            valueType:  NAME,
        }, nil
    }
    if t.Kind == "unReference" {
        i := fork.Read()
        if fork.Read().Kind != "unReference" {
            return nil, errors.New("You have a SQL syntax error near:" + string(i.Raw))
        }
        tr.Next(3)
        return &syntaxTreeNode {
            name:   "identical",
            value:  t.Raw,
            valueType: NAME,
        }, nil
    }
    stn, err := valueParser(tr)
    if err != nil {
        return nil, err
    }
    value, valueType := stn.value, stn.valueType
    return &syntaxTreeNode {
        name:   "identical",
        value:  value,
        valueType:  valueType,
    }, nil
}

func valueParser(tr *TokenReader) (*syntaxTreeNode, error) {
    t := tr.Read()
    if t.Kind == "intval" {
        return &syntaxTreeNode {
            name:   "value",
            value:  calculateInt(t.Raw),
            valueType:  INT,
        }, nil
    }
    if t.Kind == "floatval" {
        return &syntaxTreeNode {
            name:   "value",
            value:  calculateFloat(t.Raw),
            valueType:  FLOAT,
        }, nil
    }
    return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}

func filedParser(tr *TokenReader) (*syntaxTreeNode, error) {
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
            return &syntaxTreeNode {
                name:   "spot",
                value:  nil,
                child:  []*syntaxTreeNode{
                     &syntaxTreeNode {
                        name:   "identcial",
                        value:  t.Raw,
                        valueType:  NAME,
                    },
                    snt,
                },
            }, nil
        }
        return &syntaxTreeNode {
            name:   "identcial",
            value:  t.Raw,
            valueType:  NAME,
        }, nil
    }
    return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}