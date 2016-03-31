package syntax

import "errors"

func stringParser(tr *TokenReader) (*SyntaxTreeNode, error) {
    
    t := tr.Read()
    if t.Kind != "reference" {
        return nil, 
        errors.New("You have a syntax error near:" + string(t.Raw))
    }
    raws := []byte{}
    for t = tr.Read();t.Kind != "reference" && !tr.Empty();t = tr.Read() {
        raws = append(raws, t.Raw...)
    }
    if t.Kind != "reference" {
        return nil, 
        errors.New("You have a syntax error near:" + string(t.Raw))
    }
    return &SyntaxTreeNode {
        Name:       "string",
        Value:      raws,
        ValueType:  STRING,
    }, nil
}