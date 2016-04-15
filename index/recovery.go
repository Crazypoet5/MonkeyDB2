package index

import (
	"encoding/json"
	"io/ioutil"
	"unsafe"

	"../memory"

	"./csbt"
	"./mem"

	"../common"
)

var RecoveryTable = make(map[uintptr]*Index)

type savedIndex struct {
	Kind      int
	Database  string
	Table     string
	Key       string
	I         uint
	Implement string
	OldPtr    uint
}

func Restore() {
	sis := make([]savedIndex, 0)
	for _, v := range IndexList {
		si := savedIndex{
			Kind:     v.Kind,
			Database: v.Database,
			Table:    v.Table,
			Key:      v.Key,
			OldPtr:   uint(uintptr(unsafe.Pointer(v))),
		}
		if csb, ok := v.I.(*csbt.DCSBT); ok {
			si.Implement = "dcsbt"
			si.I = uint(csb.MB.RawPtr)
		}
		//TODO:Hash

		sis = append(sis, si)
	}
	data, _ := json.Marshal(sis)
	ioutil.WriteFile(common.COMMON_DIR+"\\index.json", data, 0666)
}

//Must run memory.Recovery before
func Recovery() {
	var sis []map[string]interface{}
	data, _ := ioutil.ReadFile(common.COMMON_DIR + "\\index.json")
	json.Unmarshal(data, &sis)
	for _, v := range sis {
		i := &Index{
			Kind:     int(v["Kind"].(float64)),
			Database: v["Database"].(string),
			Table:    v["Table"].(string),
			Key:      v["Key"].(string),
		}
		if v["Implement"].(string) == "dcsbt" {
			i.I = &csbt.DCSBT{
				MB: &mem.ManagedBlock{
					DataBlock: *(memory.RecoveryTable[uintptr(uint(v["I"].(float64)))]),
				},
			}
		}
		RecoveryTable[uintptr(uint(v["OldPtr"].(float64)))] = i
		IndexList = append(IndexList, i)
	}
}
