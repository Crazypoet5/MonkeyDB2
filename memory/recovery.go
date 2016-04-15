package memory

//#include <stdlib.h>
import "C"
import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"unsafe"

	"../common"
	"../log"
)

const (
	MAX_BAK_CHAN_SIZE = 10
)

var startBackup = make(chan bool, MAX_BAK_CHAN_SIZE)

var RecoveryTable = make(map[uintptr]*DataBlock)

func (b *DataBlock) SyncToFile() error {
	data, err := b.Read(0, b.Size)
	if err != nil {
		return err
	}
	filename, ok := ImageTable[b.RawPtr]
	if !ok {
		return NOT_FOUND_ADDRESS
	}
	log.WriteLog("sys", "Sync "+strconv.Itoa(int(b.RawPtr))+" to file.")
	return ioutil.WriteFile(filename, data, 0666)
}

func SyncAllImageToFile() {
	for l := DataBlockList.Front(); l != nil; l = l.Next() {
		b, ok := l.Value.(*DataBlock)
		if !ok {
			continue
		}
		data, _ := b.Read(0, b.Size)
		filename, ok := ImageTable[b.RawPtr]
		if !ok {
			continue
		}
		ioutil.WriteFile(filename, data, 0666)
	}
}

func SaveImageTable() {
	tempTable := make(map[string]string)
	for k, v := range ImageTable {
		tempTable[strconv.Itoa(int(k))] = v
	}
	data, _ := json.Marshal(tempTable)
	ioutil.WriteFile(common.COMMON_DIR+"\\image\\imageTable.json", data, 0666)
	ioutil.WriteFile(common.COMMON_DIR+"\\image\\count", []byte(strconv.Itoa(count)), 0666)
	log.WriteLog("sys", "Save image table to file.")
}

func SignalBackup() {
	startBackup <- true
}

func BackupRoutine() {
	for {
		<-startBackup
		SaveImageTable()
		SyncAllImageToFile()
	}
}

func init() {
	//	Recovery()
	//	go BackupRoutine()
}

func Restore() {
	SaveImageTable()
	SyncAllImageToFile()
}

func LoadImage(filename string) *DataBlock {
	data, _ := ioutil.ReadFile(filename)
	size := len(data)
	ip := &DataBlock{
		RawPtr: uintptr(C.malloc(C.size_t(size))),
		Size:   uint(size),
	}
	var header reflect.SliceHeader
	header.Data = uintptr(ip.RawPtr)
	header.Len = size
	header.Cap = size
	d := *(*[]byte)(unsafe.Pointer(&header))
	copy(d, data)
	ImageTable[ip.RawPtr] = filename
	DataBlockTable[ip.RawPtr] = ip
	DataBlockList.PushBack(ip)
	return ip
}

func Recovery() {
	countS, err := ioutil.ReadFile(common.COMMON_DIR + "\\image\\count")
	if err != nil {
		log.WriteLog("sys", "Recovery abort:"+err.Error())
		return
	}
	count, _ = strconv.Atoi(string(countS))
	data, err := ioutil.ReadFile(common.COMMON_DIR + "\\image\\imageTable.json")
	if err != nil {
		log.WriteLog("sys", "Recovery abort:"+err.Error())
		return
	}
	var tempTable map[string]interface{}
	json.Unmarshal(data, &tempTable)
	for k, v := range tempTable {
		ipOld, _ := strconv.Atoi(k)
		ipNew := LoadImage(v.(string))
		RecoveryTable[uintptr(ipOld)] = ipNew
	}
	log.WriteLog("sys", "Recovery "+strconv.Itoa(len(RecoveryTable))+" image files.")
}

func ShutDown() {
	SaveImageTable()
	SyncAllImageToFile()
	log.WriteLog("sys", "memory manager system shutdown.")
}
