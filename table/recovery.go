package table

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"unsafe"

	"../common"
	"../index"
	"../log"
	"../memory"
)

type SavedField struct {
	Name      string
	FixedSize bool
	Index     uint
	Size      int
	Type      int
}

type SavedTable struct {
	Name      string
	Fields    []SavedField
	FirstPage uint
	LastPage  uint
	Primary   int
}

func syncTablesToFile() {
	svts := make([]SavedTable, 0)
	for _, v := range TableList {
		svt := SavedTable{
			Name:      v.Name,
			FirstPage: uint(v.FirstPage.RawPtr),
			LastPage:  uint(v.LastPage.RawPtr),
			Primary:   v.Primary,
		}
		for _, f := range v.Fields {
			sf := SavedField{
				Name:      f.Name,
				FixedSize: f.FixedSize,
				Index:     uint(uintptr(unsafe.Pointer(f.Index))),
				Size:      f.Size,
				Type:      f.Type,
			}
			svt.Fields = append(svt.Fields, sf)
		}
		svts = append(svts, svt)
	}
	data, _ := json.Marshal(svts)
	ioutil.WriteFile(common.COMMON_DIR+"\\tables.json", data, 0666)
}

func Restore() {
	syncTablesToFile()
}

// Before run this function please run memory.RecoveryFirst
func Recovery() {
	TableList = []*Table{} //WARNING
	data, err := ioutil.ReadFile(common.COMMON_DIR + "\\tables.json")
	if err != nil {
		log.WriteLog("sys", "Tables recovery interrupt with an error :"+err.Error())
		return
	}
	var svt []map[string]interface{}
	err = json.Unmarshal(data, &svt)
	if err != nil {
		log.WriteLog("sys", "Tables recovery interrupt with an error :"+err.Error())
		return
	}
	for _, v := range svt {
		t := &Table{
			Name:    v["Name"].(string),
			Primary: int(v["Primary"].(float64)),
		}
		for _, f := range v["Fields"].([]interface{}) {
			fd := f.(map[string]interface{})
			field := Field{
				Name:      fd["Name"].(string),
				FixedSize: fd["FixedSize"].(bool),
				Size:      int(fd["Size"].(float64)),
				Type:      int(fd["Type"].(float64)),
			}
			ind := uintptr(uint(fd["Index"].(float64)))
			field.Index = index.RecoveryTable[ind]
			t.Fields = append(t.Fields, field)
		}
		firstPagePtr := uintptr(uint(v["FirstPage"].(float64)))
		firstPage := &Page{
			DataBlock: *(memory.RecoveryTable[firstPagePtr]),
		}
		t.FirstPage = firstPage
		t.LastPage = firstPage.Recovery(t, nil)
		TableList = append(TableList, t)
	}
	log.WriteLog("sys", "Recovery "+strconv.Itoa(len(TableList))+" Tables.")
}
