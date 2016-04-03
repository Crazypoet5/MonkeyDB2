package index

import (
	"./csbt"
	"./msg"
)

const (
	PRIMARY = iota
	UNIQUE
)

type Indexer interface {
	Select(key uint32, n int) []uintptr
	Insert(key uint32, value uintptr) *msg.Msg
	//InsertBat(key, value []uint, base []uintptr) msg.Msg
	Recovery()
}

type Index struct {
	Kind     int
	Database string
	Table    string
	Key      string
	I        Indexer
}

var IndexTable map[uintptr]*Index

func CreateIndex(kind int, database string, table string, key string) *Index {
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
	return i
}
