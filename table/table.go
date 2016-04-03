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
	// Skip		uint	//in every row there is a uint skip to indicate this row was deleted and how many to skip to next record
	Name      string
	FixedSize bool
	Index     *index.Index
	Size      int
	Type      int
}

type Table struct {
	Name      string
	Fields    []Field
	FirstPage *Page
}

type Index struct {
	Name   string
	Fields []Field
}
