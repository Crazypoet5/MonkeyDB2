package table

import (
	"errors"

	"../index"
)

var TableList []*Table

func CreateTable(tableName string) *Table {
	ret := &Table{
		Name:    tableName,
		Primary: -1,
	}
	ret.FirstPage = ret.NewPage()
	ret.LastPage = ret.FirstPage
	TableList = append(TableList, ret)
	return ret
}

func (t *Table) AddFiled(fieldName string, fixedSize bool, size int, ttype int, keyType int) error {
	if t.Primary != -1 && keyType == index.PRIMARY {
		return errors.New("Table has a primary key already")
	}
	if keyType == index.PRIMARY {
		t.Primary = len(t.Fields)
	}
	ind := index.CreateIndex(keyType, "db", t.Name, fieldName+"_key")
	filed := Field{
		Name:      fieldName,
		FixedSize: fixedSize,
		Size:      size,
		Type:      ttype,
		Index:     ind,
	}
	t.Fields = append(t.Fields, filed)
	return nil
}
