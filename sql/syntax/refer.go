package syntax

import "errors"

func referParser(tr *TokenReader) (*SyntaxTreeNode, error) {
	t := tr.Read()
	if t.Kind != "unreference" {
		return nil,
			errors.New("You have a syntax error near:" + string(t.Raw))
	}
	raws := []byte{}
	for t = tr.Read(); t.Kind != "unrefrence" && !tr.Empty(); {
		raws = append(raws, t.Raw...)
	}
	if t.Kind != "unreference" {
		return nil,
			errors.New("You have a syntax error near:" + string(t.Raw))
	}
	return &SyntaxTreeNode{
		Name:      "refer",
		Value:     raws,
		ValueType: NAME,
	}, nil
}
