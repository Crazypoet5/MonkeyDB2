package plan

import (
	"errors"
	"unsafe"

	"../exe"
	"../sql/syntax"
	"../table"
)

func getKVPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "getkv" {
		return nil, nil, errors.New("expected getkv but get:" + stn.Name)
	}
	tableName := string(stn.Value.([]byte))
	t := table.GetTableByName(tableName)
	if t == nil {
		return nil, nil, errors.New("table not exist")
	}
	var k []byte
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
	v, err := t.KVGetValue(k)
	if err != nil {
		return nil, nil, err
	}
	rel := exe.NewRelation()
	rel.SetColumnNames([]string{"value"})
	rel.AddRow(exe.Row{exe.NewValue(t.Fields[1].Type, v)})
	res.SetResult(1)
	return rel, res, nil
}
