package table

import (
	"../index"
)

const (
	FIELD_TYPE_INT = iota
	FIELD_TYPE_FLOAT
	FIELD_TYPE_VARCHAR
)

type Field struct {
	Name      string
	FixedSize bool
	Index     *index.Index
	Size      int
	Type      int
}

type savedRow struct {
	skip uint   //in every row there is a uint skip to indicate this row was deleted and how many to skip to next record
	size uint32 //if not fixed
	data []byte
}

type Table struct {
	Name      string
	Fields    []Field
	FirstPage *Page
	LastPage  *Page
}

type Index struct {
	Name   string
	Fields []Field
}

func GetTableByName(name string) *Table {
	for _, v := range TableList {
		if v.Name == name {
			return v
		}
	}
	return nil
}
