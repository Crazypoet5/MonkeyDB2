package network

import (
	"fmt"
	"strconv"
	"unsafe"

	"../exe"
	"../plan"
)

type Pack struct {
	Head uint32
	Len  uint32
	Type uint32
	Data []byte
}

const (
	CREATE_CONNECTION = iota
	CLOSE_CONNECTION
	DIRECT_QUERY
	RESPONSE
)

func Encode(rel *exe.Relation, res *plan.Result) *Pack {
	relJson := "["
	if rel != nil && len(rel.Rows) != 0 {
		for r := 0; r < len(rel.Rows); r++ {
			if r != 0 {
				relJson += ","
			}
			relJson += "{"
			for i := 0; i < len(rel.Rows[r]); i++ {
				if i != 0 {
					relJson += ","
				}
				if i < len(rel.ColumnNames) {
					relJson += "\"" + rel.ColumnNames[i] + "\":"
				} else {
					relJson += "\"" + strconv.Itoa(i) + "\":"
				}
				switch rel.Rows[r][i].Kind {
				case exe.INT:
					i := *(*int)((unsafe.Pointer)(&rel.Rows[r][i].Raw[0]))
					relJson += strconv.Itoa(i)
				case exe.FLOAT:
					f := *(*float64)((unsafe.Pointer)(&rel.Rows[r][i].Raw[0]))
					relJson += strconv.FormatFloat(f, 'f', -1, 64)
				default: //exe.STRING
					relJson += "\"" + string(rel.Rows[r][i].Raw) + "\""
				}
			}
			relJson += "}"
		}
	}
	relJson += "]"
	resJson := "{\"affectedRows\":%d,\"usedTime\":%d}"
	resJson = fmt.Sprintf(resJson, res.AffectedRows, res.UsedTime)
	json := []byte(fmt.Sprintf("{\"relation\":%s,\"result\":%s}", relJson, resJson))
	return &Pack{
		Head: 2016,
		Len:  uint32(len(json) + 12),
		Type: RESPONSE,
		Data: json,
	}
}
