package mem

import (
    "../../memory"
    "../../log"
)

const (
    MAX_POWER       =    13
    NORMAL_SIZE     =    24 + 16 * 1024 * 32
    NORMAL_SIZE_2X  =    2 * NORMAL_SIZE
)

type ManagedBlock struct {
    memory.DataBlock
}

type header struct {
    root        uint
    min         uint
    max         uint
    freeList    [MAX_POWER]uint
}

type freeListElement struct {
    next        uint
    data        interface{}
}

func NewManagedBlockWithSize(size int) *ManagedBlock {
    db, err := memory.CreateImage(size)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    mb := &ManagedBlock {   //This struct actually save 8 byte pointer only
        DataBlock:  *db,
    }
    mb.Init()
    return mb
}

func NewManagedBlock() *ManagedBlock {
    return NewManagedBlockWithSize(NORMAL_SIZE)
}

func (mb *ManagedBlock) ListGetNext(e uint) uint {
    data, err := mb.Read(e, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) ListGetBack(list uint) uint {
    ret := list
    for e := list;e != 0;e = mb.ListGetNext(e) {
        ret = e
    }
    return ret
}

func (mb *ManagedBlock) ListGetLength(n int) int {
    i := 0
    for e := mb.GetFreeList(n);e != 0;e = mb.ListGetNext(e) {
        i++
    }
    return i
}

func (mb *ManagedBlock) ListPushBack(n int, e uint) {
    list := mb.GetFreeList(n)
    back := mb.ListGetBack(list)
    if back == 0 {
        mb.Write(FREE_LIST_OFFSET + 8 * uint(n), uint2bytes(e))
    } else {
        mb.Write(back, uint2bytes(e))
    }
    mb.Write(e, uint2bytes(0))
}

func (mb *ManagedBlock) ListPopBack(n int) uint {
    list := mb.GetFreeList(n)
    next := list
    p := next
    if p == 0 {
        return 0
    }
    i := 0
    next =  mb.ListGetNext(next)
    if next == 0 {  // Only one element       
        mb.Write(FREE_LIST_OFFSET + 8 * uint(n), uint2bytes(0))
        return list
    }
    for ;next != 0;next = mb.ListGetNext(next) {
        if i > 0 {
            p = mb.ListGetNext(p)
        }
        i++
    }
    ret := mb.ListGetNext(p)
    mb.Write(p, uint2bytes(0))
    return ret
}

func (mb *ManagedBlock) ListPopFront(n int) uint {
    front := mb.GetFreeList(n)
    if front == 0 {
        return 0
    }
    next := mb.ListGetNext(front)
    mb.Write(FREE_LIST_OFFSET + 8 * uint(n), uint2bytes(next))
    return front
}

func (mb *ManagedBlock) ListPushFront(n int, e uint) {
    front := mb.GetFreeList(n)
    if front == 0 {
        mb.Write(e, uint2bytes(0))
        mb.Write(FREE_LIST_OFFSET + 8 * uint(n), uint2bytes(e))
        return
    }
    mb.Write(e, uint2bytes(front))
    mb.Write(FREE_LIST_OFFSET + 8 * uint(n), uint2bytes(e))
    return
}

func (mb *ManagedBlock) Init() {
    used := uint(24 + 8 * MAX_POWER)
    for {
        for i := 0;i < MAX_POWER;i++ {
            for n := 0;n < MAX_POWER - i;n++ {
                if mb.DataBlock.Size - used < 64 * uint(i + 1) {
                    return
                }
                mb.ListPushFront(i, used)
                used += 64 * uint(i + 1)
            }
        }
    }
}
