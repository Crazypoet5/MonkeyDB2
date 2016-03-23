package mem

import (
    "../../log"
)

type Node struct {
    IsLeaf      byte
    Reversed    byte
    KeyNum      uint16
    Key         [13]uint32
    Child       uint
}

func (mb *ManagedBlock) InitNode(p uint) {
    mb.db.Write(p, make([]byte, 64))
}

func (mb *ManagedBlock) NewNodes(n int) uint {
    nodes := mb.Malloc(n)
    for i := 0;i < n;i++ {
        mb.InitNode(nodes + uint(i) * 64)
    }
    return nodes
}

func (mb *ManagedBlock) SetNodeKey(node uint, index int, key uint32) {
    mb.db.Write(node + 4 + uint(index) * 4, uint322bytes(key))
}

func (mb *ManagedBlock) GetNodeKey(node uint, index int) uint32 {
    data, err := mb.db.Read(node + 4 + uint(index) * 4, 4)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint32(data)
}

func (mb *ManagedBlock) GetNodeKeyNum(node uint) int {
    data, err := mb.db.Read(uint(node + 2), 2)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return int(bytes2uint16(data))
}

func (mb *ManagedBlock) SetNodeKeyNum(node uint, keyNum int) {
    data := uint162bytes(uint16(keyNum))
    mb.db.Write(uint(node + 2), data)
}

func (mb *ManagedBlock) GetChild(node uint, index int) uint {
    data, err := mb.db.Read(uint(node + 56), 8)
    if err != nil {
        log.WriteLog("err", err.Error())
    }
    return bytes2uint(data)
}