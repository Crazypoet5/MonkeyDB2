package table

type Field struct {
    Name        string
    FixedSize   bool
}

type Table struct {
    Name        string
    Fields      []Field
    Indexs      []Index
}

type Page struct {
    TableHash   uint
    Len         uint
    Cap         uint
    Base        uintptr   
}

type Index struct {
    Unique      bool
    Name        string
    Fields      []Field
}