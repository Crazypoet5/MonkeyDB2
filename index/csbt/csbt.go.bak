//      DCSBT Index @InsZVA 2015
//      Based on WANG Sheng, QIN Xiaolin, SHEN Yao, et al. Research on durable CSB+-tree indexing technology. Journal of Frontiers of Computer Science and Technology, 2015, 9(2): 182-192.

package csbt

import (
    "unsafe"
    "../../mempool"
)

const (
    CACHE_LINE = 64
    NODE_NUM = (CACHE_LINE - 24) / 8
    ALLOCATE = 16 * 1024
    LEAF_SIZE = (NODE_NUM / 3) * 24 + 24
    NODE_SIZE = NODE_NUM * 8 + 24
)

//interface{} with a 8 byte data is sized with 16 bytes in GOlang

type Node struct {
    IsLeaf      byte
    Reversed    byte
    KeyNum      uint16
    Key         [13]uint32
    Child       uintptr
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

type IndexHeader struct {
    root    interface{}
    min,max *CSBTLeaf
}

type DCSBT struct {
    indexHeader *IndexHeader
    BaseAddr uintptr
    freePos     uintptr
}

type CSBTNode struct {
    keyNum int
    child interface{}
    key [NODE_NUM]uint
}

type CSBTLeaf struct {
    keyNum int
    key [NODE_NUM / 3]uint
    data [NODE_NUM / 3]uint
    base [NODE_NUM / 3]uintptr
    left, right *CSBTLeaf
}

func newCDBTLeaf(freePos uintptr) *CSBTLeaf {
    r := makeBytes(freePos, LEAF_SIZE)
    root := (*CSBTLeaf)(unsafe.Pointer(&r[0]))
    return root
}

func NewDCSBT() *DCSBT {
    malloc := mempool.Malloc(ALLOCATE)
    freePos := uintptr(unsafe.Pointer(&malloc[0]))
    baseAddr := freePos
    bytes := makeBytes(freePos, 32)
    indexHeader := (*IndexHeader)(unsafe.Pointer(&bytes[0]))
    freePos += 32
    root := newCDBTLeaf(freePos)
    freePos += LEAF_SIZE
    if unsafe.Sizeof(*root) != LEAF_SIZE {
        panic(int(unsafe.Sizeof(*root)))
    }
    indexHeader.root = root
    return &DCSBT {
        indexHeader:    indexHeader,
        BaseAddr:       baseAddr,
        freePos:        freePos,
    }
}
