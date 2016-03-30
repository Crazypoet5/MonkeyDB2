package syntax

import (
    "errors"
)

func typeParser(tr *TokenReader) (*syntaxTreeNode, error) {
    t := tr.Read()
    if t.Kind == "types" {
        return &syntaxTreeNode {
            name:       "type",
            value:      t.Raw,
            valueType:  NAME,
        }, nil
    }
    if t.Kind == "floatval" {
        return &syntaxTreeNode {
            name:       "value",
            value:      t.Raw,
            valueType:  NAME,
        }, nil
    }
    return nil, errors.New("You have a syntax error near:" + string(t.Raw))
}

func columndefineParser(tr *TokenReader) (*syntaxTreeNode, error) {
    t := tr.Read()
    if t.Kind != "identical" {
        return nil, errors.New("You have a syntax error near:" + string(t.Raw))
    }
    tpNode, err := typeParser(tr)
    if err != nil {
        return nil, err
    }
    first := &syntaxTreeNode {
        name:       "ColumnDefine",
        value:      nil,
        child:      []*syntaxTreeNode {
            &syntaxTreeNode {
                name:   "identical",
                value:  t.Raw,
                valueType:  NAME,
            },
            tpNode,
    }}
    fork := tr.Fork()
    t2 := fork.Read()
    if t2.Kind == "structs" && string(t2.Raw) == "," {
        tr.Next(1)
        nexts, err := columndefineParser(tr)
        if err != nil {
            return nil, err
        }
        return &syntaxTreeNode {
            name:       "dot",
            child:      []*syntaxTreeNode {
                first, nexts,
            },
        }, nil
    }
    return first, nil    
}

func createtableParser(tr *TokenReader) (*syntaxTreeNode, error) {
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
    return &syntaxTreeNode {
        name:       "createtable",
        child:      []*syntaxTreeNode {
            &syntaxTreeNode {
                name:       "identical",
                value:      t3.Raw,
                valueType:  NAME,
            },
            columndefine,
        },
    }, nil
}