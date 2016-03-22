package mem

import (
    "../../log"
)

const (
    // These constant shows where the indexheader was storaged
    ROOT_POINTER_OFFSET     =       0
    MIN_POINTER_OFFSET      =       8
    MAX_POINTER_OFFSET      =       16
)

func (mb *ManagedBlock) GetRoot() uint {
    data, err := mb.db.Read(ROOT_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetRoot(r uint) {
    data := uint2bytes(r)
    _, err := mb.db.Write(ROOT_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMin() uint {
    data, err := mb.db.Read(MIN_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMin(p uint) {
    data := uint2bytes(p)
    _, err := mb.db.Write(MIN_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMax() uint {
    data, err := mb.db.Read(MAX_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMax(p uint) {
    data := uint2bytes(p)
    _, err := mb.db.Write(MAX_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}
