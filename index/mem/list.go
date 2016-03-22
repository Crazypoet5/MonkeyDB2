package mem

import (
    "container/list"
    "../../memory"
    "../../log"
)

const (
    MAX_POWER   =    10
)

type ManagedBlock struct {
    db  *memory.DataBlock
    freeList [MAX_POWER]list.List
    used   int
}

func NewManagedBlock(size int) *ManagedBlock {
    db, err := memory.CreateImage(size)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    mb := &ManagedBlock {
        db:         db,
    }
    mb.Init()
    return mb
}

func (mb *ManagedBlock) Init() {
    free := mb.db.Size - 24     // Expect header
    
}