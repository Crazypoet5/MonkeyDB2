package mem

import (
    "../../log"
)

const (
    // These constant shows where the indexheader was storaged
    ROOT_POINTER_OFFSET     =       0
    MIN_POINTER_OFFSET      =       8
    MAX_POINTER_OFFSET      =       16
    FREE_LIST_OFFSET        =       24
)

func (mb *ManagedBlock) GetRoot() uint {
    data, err := mb.Read(ROOT_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetRoot(r uint) {
    data := uint2bytes(r)
    _, err := mb.Write(ROOT_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMin() uint {
    data, err := mb.Read(MIN_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMin(p uint) {
    data := uint2bytes(p)
    _, err := mb.Write(MIN_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMax() uint {
    data, err := mb.Read(MAX_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMax(p uint) {
    data := uint2bytes(p)
    _, err := mb.Write(MAX_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetFreeList(i int) uint {
    data, err := mb.Read(FREE_LIST_OFFSET + 8 * uint(i), 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}