package table

import (
    "../index"
)

const (
    FIELD_TYPE_INT  =   iota
    FIELD_TYPE_FLOAT
    FIELD_TYPE_VARCHAR
)

type Field struct {
    Name        string
    FixedSize   bool
    Index       index.Index
    Size        int
    Type        int
}

type Table struct {
    Name        string
    Fields      []Field
    Indexs      []Index
    FirstPage   *Page 
}

type Index struct {
    Name        string
    Fields      []Field
}