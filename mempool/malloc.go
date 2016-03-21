package mempool

import (
    "../log"
    "reflect"
    "unsafe"
    "../memory"
)

var UsageTable = make(map[uintptr]*memory.DataBlock)

// Please use GetFree instead or you might make an error
func Malloc(size int) []byte{
    db, _ := memory.CreateImage(size)
    UsageTable[db.RawPtr] = db
    var header reflect.SliceHeader
    header.Data = db.RawPtr
    header.Len = size
    header.Cap = size
    b := *(*[]byte)(unsafe.Pointer(&header))
    for i := 0;i < size;i++ {
        b[i] = 0
    }
    return b
}

// Please use Release instead or you might make an error
func Free(p []byte) {
    header := (*reflect.SliceHeader)(unsafe.Pointer(&p))
    db, ok := UsageTable[header.Data]
    if !ok {
        log.WriteLog("err", "Free an unregistered data block.")
        return
    }
    memory.ReleaseImage(db)
}
