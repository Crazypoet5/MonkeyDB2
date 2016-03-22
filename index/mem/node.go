package mem

import (
    
)

type Node struct {
    IsLeaf      byte
    Reversed    byte
    KeyNum      uint16
    Key         [13]uint32
    Child       uint
}

type Leaf struct {
    IsLeaf      byte
    KeyNum      byte
    Key         [3]uint32
    Value       [3]uintptr
    Left,Right  uintptr
    Reversed    uint32
    Reversed2   uint16
}

func (mb *ManagedBlock) InitNode(p uint) {
    mb.db.Write(p, []byte{1})
    mb.db.Write(p + 1, []byte{0})
    mb.db.Write(p + 2, uint162bytes(0))
    for i := 0;i < MAX_POWER;i++ {
        mb.db.Write(p + 4 * uint(i + 1), uint322bytes(0))
    }
    mb.db.Write(p + 56, uint2bytes(0))
}

func (mb *ManagedBlock) NewNodes(n int) uint {
    nodes := mb.Malloc(n)
    for i := 0;i < n;i++ {
        mb.InitNode(nodes + uint(i) * 64)
    }
    return nodes
}