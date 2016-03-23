package mem

import (
    "../../log"
)

type leaf struct {
    isLeaf      byte
    keyNum      byte
    reserved    uint16
    key         [3]uint32
    value       [3]uintptr
    left,right  uint
    reserved2   uint
}

func (mb *ManagedBlock) InitLeaf(p uint) {
    mb.db.Write(p, make([]byte, 64))
}

func (mb *ManagedBlock) IsLeaf(p uint) bool {
    data, err := mb.db.Read(p, 1)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return data[0] >= 0
}

func (mb *ManagedBlock) NewLeaves(n int) uint {
    leaves := mb.Malloc(n)
    for i := 0;i < n;i++ {
        mb.InitLeaf(leaves + uint(i) * 64)
    }
    return leaves
}

func (mb *ManagedBlock) GetLeafKeyNum(leaf uint) int {
    data, err := mb.db.Read(leaf + 1, 1)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return int(data[0])
}

func (mb *ManagedBlock) SetLeafKeyNum(leaf uint, keyNum int) {
    mb.db.Write(leaf + 1, []byte{byte(keyNum)})
}

func (mb *ManagedBlock) GetLeafKey(leaf uint, index int) uint32 {
    data, err := mb.db.Read(leaf + 4 + uint(index) * 4, 4)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint32(data)
}

func (mb *ManagedBlock) SetLeafKey(leaf uint, index int, key uint32) {
    mb.db.Write(leaf + 4 + uint(index) * 4, uint322bytes(key))
}

func (mb *ManagedBlock) GetLeafValue(leaf uint, index int) uintptr {
    data, err := mb.db.Read(leaf + 16 + uint(index) * 8, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return uintptr(bytes2uint(data))
}

func (mb *ManagedBlock) SetLeafValue(leaf uint, index int, value uintptr) {
    mb.db.Write(leaf + 16 + uint(index) * 8, uint2bytes(uint(value)))
}

func (mb *ManagedBlock) GetLeafLeft(leaf uint) uint {
    data, err := mb.db.Read(leaf + 40, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetLeafLeft(leaf, left uint) {
    mb.db.Write(leaf + 40, uint2bytes(uint(left)))
}

func (mb *ManagedBlock) GetLeafRight(leaf uint) uint {
    data, err := mb.db.Read(leaf + 48, 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}

func (mb *ManagedBlock) SetLeafRight(leaf, right uint) {
    mb.db.Write(leaf + 48, uint2bytes(uint(right)))
}
