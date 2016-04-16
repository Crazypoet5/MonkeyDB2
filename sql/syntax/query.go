package syntax

import (
	"errors"
)

func Parser(tr *TokenReader) (*SyntaxTreeNode, error) {
	fork := tr.Fork()
	t := fork.Read()
	if t.Kind == "keyword" && string(t.Raw) == "create" {
		return createtableParser(tr)
	}
	if t.Kind == "keyword" && string(t.Raw) == "insert" {
		return insertParser(tr)
	}
	if t.Kind == "keyword" && string(t.Raw) == "dump" {
		return dumpParser(tr)
	}
	if t.Kind == "keyword" && string(t.Raw) == "select" {
		return selectParser(tr)
	}
	return nil, errors.New("Unsupported syntax!")
}
