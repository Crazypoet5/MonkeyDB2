package table

var TableList []*Table

func CreateTable(tableName string) *Table {
    return &Table {
        Name:       tableName,
        FirstPage:  NewPage(),
    }
}

func (t *Table) AddFiled(fieldName string, fixedSize bool, size int, ttype int) {
    filed := Field {
        Name:       fieldName,
        FixedSize:  fixedSize,
        Size:       size,
        Type:       ttype,
    }
    t.Fields = append(t.Fields, filed)
}