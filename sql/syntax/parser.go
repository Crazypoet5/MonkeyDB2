package syntax

import (
    "fmt"
)

/*
    createquery = create table identical ( columndefine )
    columndefine = identical type | identical type , columndefine
    type = int | float
    query = selectquery | insertquery | updatequery | deletequery | createquery
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
    if stn == nil {
        return
    }
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
