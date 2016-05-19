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
				if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
					return nil, errors.New("Unexpect end:" + string(t.Raw))
				}
				return nil, err
			}
			return stn, nil
		}
		tr.pos = fork2.pos
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "createkv" {
		stn, err := createKVParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "set" {
		stn, err := setKVParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "get" {
		stn, err := getParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "remove" {
		stn, err := removeParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "insert" {
		stn, err := insertParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "dump" {
		stn, err := dumpParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "select" {
		stn, err := selectParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "delete" {
		stn, err := deleteParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "drop" {
		fork2 := tr.Fork()
		stn, err := dropParser(fork2)
		if err != nil {
			stn, err := dropIndexParser(tr)
			if err != nil {
				return nil, err
			}
			if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
				return nil, errors.New("Unexpect end:" + string(t.Raw))
			}
			return stn, nil
		}
		tr.pos = fork2.pos
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "update" {
		stn, err := updateParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	if t.Kind == "keyword" && string(t.Raw) == "show" {
		stn, err := showParser(tr)
		if err != nil {
			return nil, err
		}
		if t := tr.Read(); t.Kind != "" && string(t.Raw) != ";" {
			return nil, errors.New("Unexpect end:" + string(t.Raw))
		}
		return stn, nil
	}
	return nil, errors.New("Unsupported syntax!")
}
