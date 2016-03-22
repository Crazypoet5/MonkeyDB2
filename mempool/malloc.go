package mempool

import (
    "../log"
    "../memory"
)

// Please use GetFree instead or you might make an error
func Malloc(size int) *memory.DataBlock{
    db, err := memory.CreateImage(size)
    if err != nil {
        log.WriteLog("err", "Malloc with an error : " + err.Error())
    }
    return db
}

// Please use Release instead or you might make an error
func Free(db *memory.DataBlock) {
    err := memory.ReleaseImage(db)
    if err != nil {
        log.WriteLog("err", "Malloc with an error : " + err.Error())
    }
}
