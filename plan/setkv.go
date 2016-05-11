package plan

import (
	"errors"
	"unsafe"

	"../exe"
	"../sql/syntax"
	"../table"
)

func setKVPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "setkv" {
		return nil, nil, errors.New("expected setkv but get:" + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("table not exist")
	}
	var k, v []byte
	switch stn.Child[0].ValueType {
	case syntax.INT:
		k = make([]byte, 8)
		i := stn.Child[0].Value.(int)
		p_bytes := *(*[8]byte)(unsafe.Pointer(&i))
		for t := 0; t < 8; t++ {
			k[t] = p_bytes[t]
		}
	case syntax.FLOAT:
		k = make([]byte, 8)
		i := stn.Child[0].Value.(float64)
		p_bytes := *(*[8]byte)(unsafe.Pointer(&i))
		for t := 0; t < 8; t++ {
			k[t] = p_bytes[t]
		}
	default:
		k = stn.Child[0].Value.([]byte)
	}
	switch stn.Child[1].ValueType {
	case syntax.INT:
		v = make([]byte, 8)
		i := stn.Child[1].Value.(int)
		p_bytes := *(*[8]byte)(unsafe.Pointer(&i))
		for t := 0; t < 8; t++ {
			v[t] = p_bytes[t]
		}
	case syntax.FLOAT:
		v = make([]byte, 8)
		i := stn.Child[1].Value.(float64)
		p_bytes := *(*[8]byte)(unsafe.Pointer(&i))
		for t := 0; t < 8; t++ {
			v[t] = p_bytes[t]
		}
	default:
		v = stn.Child[1].Value.([]byte)
	}

	err := t.KVSetValue(k, v)
	if err != nil {
		return nil, nil, err
	}
	res.SetResult(1)
	return nil, res, nil
}
