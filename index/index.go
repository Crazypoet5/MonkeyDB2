package index

import (
    "./csbt"
)

const (
    CSBT    =   iota
)

type Indexer interface {
    Select(key uint, n int) []uint
    Insert(key, value uint, base uintptr) error 
    InsertBat(key, value []uint, base []uintptr) error
    Recovery(oldPtr uintptr) error
}

type Index struct {
    Kind        int         
    Database    string
    Table       string
    Key         string
    I           Indexer
}

var IndexTable map[uintptr]Index

func CreateIndex(kind int, database string, table string, key string) *Index {
    i := &Index {
        Kind:       kind,
        Database:   database,
        Table:      table,
        Key:        key, 
    }
    switch i.Kind {
        case CSBT:
        i.I = csbt.NewDCSBT()
    }
    return i
}