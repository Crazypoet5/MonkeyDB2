package table

import (
	"errors"

	"../exe"
	"../index"
	"../memory"
)

const (
	KV_INT = iota
	KV_FLOAT
	KV_STRING
)

func CreateKVTable(tableName string, keyType int, valueType int) *Table {
	ret := CreateTable(tableName)
	switch keyType {
	case KV_INT:
		ret.AddFiled("key", true, 8, exe.INT, index.PRIMARY)
	case KV_FLOAT:
		ret.AddFiled("key", true, 8, exe.FLOAT, index.PRIMARY)
	case KV_STRING:
		ret.AddFiled("key", false, 0, exe.STRING, index.PRIMARY)
	}
	switch valueType {
	case KV_INT:
		ret.AddFiled("value", true, 8, exe.INT, -1)
	case KV_FLOAT:
		ret.AddFiled("value", true, 8, exe.FLOAT, -1)
	case KV_STRING:
		ret.AddFiled("value", false, 0, exe.STRING, -1)
	}
	return ret
}

func (t *Table) KVGetValue(key []byte) ([]byte, error) {
	if t.Primary < 0 {
		return nil, errors.New("this table cannot be used as a kv table")
	}
	ind := t.Fields[t.Primary].Index
	if ind == nil {
		return nil, errors.New(("this table cannot be used as a kv table"))
	}
	cur := ind.I.Select(index.BKDRHash(key))
	if cur == nil {
		return nil, errors.New("not found key")
	}
	ptr, offset := cur.Read()
	if ptr == 0 {
		return nil, errors.New("not found key")
	}
	p := &Page{
		DataBlock: *memory.DataBlockTable[ptr],
	}
	reader := p.NewReader()
	reader.currentPtr = offset
	rel := reader.PeekRecord()
	return rel.Rows[0][1].Raw, nil
}

func (t *Table) KVSetValue(key []byte, value []byte) error {
	if t.Primary < 0 {
		return errors.New("this table cannot be used as a kv table")
	}
	ind := t.Fields[t.Primary].Index
	if ind == nil {
		return errors.New(("this table cannot be used as a kv table"))
	}
	cur := ind.I.Select(index.BKDRHash(key))
	m := make(map[int]int)
	m[0] = 0
	m[1] = 1
	data := make([][][]byte, 1)
	data[0] = append([][]byte{}, key)
	data[0] = append(data[0], value)
	if cur == nil {
		err := t.Insert(m, data)
		if err != nil {
			return err
		}
		return nil
	}
	tab, offset := cur.Read()
	if tab == 0 {
		err := t.Insert(m, data)
		if err != nil {
			return err
		}
		return nil
	}
	cur.Delete()
	page := &Page{
		DataBlock: *memory.DataBlockTable[tab],
	}
	page.DeleteFromOffset(offset)
	t.Insert(m, data)
	return nil
}

func (t *Table) KVRemove(key []byte) error {
	if t.Primary < 0 {
		return errors.New("this table cannot be used as a kv table")
	}
	ind := t.Fields[t.Primary].Index
	if ind == nil {
		return errors.New(("this table cannot be used as a kv table"))
	}
	cur := ind.I.Select(index.BKDRHash(key))
	if cur == nil {
		return errors.New("the key has already been deleted")
	}
	tab, offset := cur.Read()
	if tab == 0 {
		return errors.New("there's a internal error")
	}
	cur.Delete()
	page := &Page{
		DataBlock: *memory.DataBlockTable[tab],
	}
	page.DeleteFromOffset(offset)
	return nil
}
