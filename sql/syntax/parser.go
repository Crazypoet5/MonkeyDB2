package syntax

import (
    "errors"
    "strconv"
    "fmt"
)

/*
    query = selectquery | insertquery | updatequery | deletequery
    selectquery = select fileds [from relation] [where logical]
    logical = expression reloperation expression | expression relperation expression logoperation logical
    fileds = filed | filed, fields
    realation = expression | relation.expression
    filed = identical | identical.filed
    expression = identical | `identical` | value
    value = int | number
*/
    
const (
    NULL    =   iota
    INT
    FLOAT
    NAME
    STRING
)

type syntaxTreeNode struct {
    name    string
    child   []*syntaxTreeNode
    value   interface{}
    valueType   int
}

func (stn *syntaxTreeNode) Print(tabs int) {
    tab := ""
    for i := 0;i < tabs;i++ {
        tab += "\t"
    }
    if stn.value != nil {
        fmt.Println(tab + stn.name, ":", stn.value)
    } else {
        fmt.Println(tab + stn.name)
    }
    
    for _, c := range stn.child {
        c.Print(tabs + 1)
    }
}

type parser func(*TokenReader) *syntaxTreeNode

func calculateInt(raw []byte) int {
    i, _ := strconv.Atoi(string(raw))
    return i
}

func calculateFloat(raw []byte) float64 {
    f, _ := strconv.ParseFloat(string(raw), 64)
    return f
}

func expressionParser(tr *TokenReader) (*syntaxTreeNode, error) {
    fork := tr.Fork()
    t := fork.Read()
    if t.Kind == "identical" {
        tr.pos++
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
        tr.pos += 3
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
    if t.Kind == "int" {
        return &syntaxTreeNode {
            name:   "value",
            value:  calculateInt(t.Raw),
            valueType:  INT,
        }, nil
    }
    if t.Kind == "float" {
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
        tr.pos++
        pre := fork.Read()
        if pre.Kind == "spot" {
            tr.pos++
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