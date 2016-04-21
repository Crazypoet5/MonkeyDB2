package index

import (
	"./csbt"
	"./cursor"
)

const (
	PRIMARY = iota
	UNIQUE
	HASH //TODO
)

var IndexList []*Index

type Indexer interface {
	Select(key uint32) cursor.Cursor
	Insert(k uint32, v uintptr)
	Delete(k uint32)
	Recovery()
}

type Index struct {
	Kind     int
	Database string
	Table    string
	Key      string
	I        Indexer
}

func CreateIndex(kind int, database string, table string, key string) *Index {
	if kind == -1 {
		return nil
	}
	i := &Index{
		Kind:     kind,
		Database: database,
		Table:    table,
		Key:      key,
	}
	switch i.Kind {
	case PRIMARY:
		i.I = csbt.NewDCSBT()
	case UNIQUE:
		i.I = csbt.NewDCSBT()
	}
	IndexList = append(IndexList, i)
	return i
}
