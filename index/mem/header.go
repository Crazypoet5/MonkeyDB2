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

func uint2bytes(p uintptr) []byte {
    b := make([]byte, 8)
    for i := 0;i < 8;i++ {
         b[i] = (byte)(p)
         p >>= 8
    }
    return b
}

func bytes2uint(b []byte) uintptr {
    p := uintptr(0)
    for i := 7;i >= 0;i-- {
        p <<= 8
        p |= uintptr(b[i])
    }
    return p
}

func (mb *ManagedBlock) GetRoot() uintptr {
    data, err := mb.db.Read(ROOT_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetRoot(r uintptr) {
    data := uint2bytes(r)
    _, err := mb.db.Write(ROOT_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMin() uintptr {
    data, err := mb.db.Read(MIN_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMin(p uintptr) {
    data := uint2bytes(p)
    _, err := mb.db.Write(MIN_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}

func (mb *ManagedBlock) GetMax() uintptr {
    data, err := mb.db.Read(MAX_POINTER_OFFSET, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetMax(p uintptr) {
    data := uint2bytes(p)
    _, err := mb.db.Write(MAX_POINTER_OFFSET, data)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
}
