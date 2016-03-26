package table

import (
    "../index"
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
    FirstPage   
}

type Index struct {
    Name        string
    Fields      []Field
}