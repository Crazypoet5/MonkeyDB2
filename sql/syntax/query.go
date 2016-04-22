package syntax

import (
	"errors"
)

func Parser(tr *TokenReader) (*SyntaxTreeNode, error) {
	fork := tr.Fork()
	t := fork.Read()
	if t.Kind == "keyword" && string(t.Raw) == "create" {
		fork2 := tr.Fork()
		stn, err := createtableParser(fork2)
		if err != nil {
			stn, err := createIndexParser(tr)
			if err != nil {
				return nil, err
			}
			return stn, nil
		}
		return stn, nil
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
	if t.Kind == "keyword" && string(t.Raw) == "delete" {
		return deleteParser(tr)
	}
	if t.Kind == "keyword" && string(t.Raw) == "drop" {
		fork2 := tr.Fork()
		stn, err := dropParser(fork2)
		if err != nil {
			stn, err := dropIndexParser(tr)
			if err != nil {
				return nil, err
			}
			return stn, nil
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "update" {
		return updateParser(tr)
	}
	return nil, errors.New("Unsupported syntax!")
}
