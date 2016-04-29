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

func Encode(rel *exe.Relation, res *plan.Result, err error) *Pack {
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
					relJson += "\"" + escape(string(rel.Rows[r][i].Raw)) + "\""
				}
			}
			relJson += "}"
		}
	}
	relJson += "]"
	resJson := "{\"affectedRows\":%d,\"usedTime\":%d}"
	if res == nil {
		resJson = fmt.Sprintf(resJson, 0, 0)
	} else {
		resJson = fmt.Sprintf(resJson, res.AffectedRows, res.UsedTime)
	}

	errStr := "null"
	if err != nil {
		errStr = "\"" + escape(err.Error()) + "\""
	}
	json := []byte(fmt.Sprintf("{\"relation\":%s,\"result\":%s,\"error\":%s}", relJson, resJson, errStr))
	return &Pack{
		Head: 2016,
		Len:  uint32(len(json) + 12),
		Type: RESPONSE,
		Data: json,
	}
}

func MakePack(tp uint32, data []byte) *Pack {
	return &Pack{
		Head: 2016,
		Len:  uint32(len(data) + 12),
		Type: tp,
		Data: data,
	}
}

func escape(raw string) string {
	var ret string
	bytes := []byte(raw)
	for i := 0; i < len(bytes); i++ {
		switch bytes[i] {
		case 10:
			ret += "\\n"
		case 13:
			ret += "\\r"
		case '\t':
			ret += "\\t"
		case '"':
			ret += "\""
		default:
			ret += string(bytes[i])
		}
	}
	return ret
}
