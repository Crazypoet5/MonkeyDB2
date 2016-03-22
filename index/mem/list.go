package mem

import (
    "container/list"
    "../../memory"
    "../../log"
)

const (
    MAX_POWER       =    13
    NORMAL_SIZE     =    24 + 16 * 1024 * 32
    NORMAL_SIZE_2X  =    2 * NORMAL_SIZE
)

type ManagedBlock struct {
    db  *memory.DataBlock
    freeList [MAX_POWER]list.List
    used   uint
}

func NewManagedBlockWithSize(size int) *ManagedBlock {
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

func NewManagedBlock() *ManagedBlock {
    return NewManagedBlockWithSize(NORMAL_SIZE)
}

func (mb *ManagedBlock) Init() {
    free := mb.db.Size - 24     // Expect header
    percent := free / MAX_POWER
    fp := uint(24)
    inited := uint(24)
    for i := uint(0);i < MAX_POWER;i++ {
        for j := uint(0);j < percent / 64 / (i + 1);j++ {
            mb.freeList[i].PushBack(fp)
            fp += uint(64 * (i + 1))
            inited += uint(64 * (i + 1))
        }
    }
    if mb.db.Size - inited > 0 {
        for i := MAX_POWER - 1;i >= 0;i-- {
            for mb.db.Size - inited > 64 * uint(i + 1) {
                mb.freeList[i].PushBack(fp)
                fp += 64 * uint(i + 1)
                inited += 64 * uint(i + 1)
            }
        }
    }
}
