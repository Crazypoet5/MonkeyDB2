//      DCSBT Index @InsZVA 2015
//      Based on WANG Sheng, QIN Xiaolin, SHEN Yao, et al. Research on durable CSB+-tree indexing technology. Journal of Frontiers of Computer Science and Technology, 2015, 9(2): 182-192.

package csbt

import (
    "unsafe"
    "errors"
    "../../mempool"
    "reflect"
)

const (
    CACHE_LINE = 64
    NODE_NUM = (CACHE_LINE - 24) / 8
    ALLOCATE = 16 * 1024
    LEAF_SIZE = (NODE_NUM / 3) * 24 + 24
)

//interface{} with a 8 byte data is sized with 16 bytes in GOlang

type IndexHeader struct {
    root    interface{}
    min,max *CSBTLeaf
}

type DCSBT struct {
    indexHeader *IndexHeader
    baseAddr uintptr
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

func makeBytes(p uintptr, n int) []byte {
    var b []byte
    header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    header.Cap = n
    header.Len = n
    header.Data = p
    return b
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
        baseAddr:       baseAddr,
        freePos:        freePos,
    }
}

func binarySearchL(target uint, array [NODE_NUM / 3]uint, l, r, w int) (int, bool) {
    if l > r {
        if l >= w {
            return l, false
        }
        return r, false
    }
    mid := (l + r) >> 1
    if array[mid] == target {
        return mid, true
    }
    if array[mid] < target {
        return binarySearchL(target, array, mid + 1, r, w)
    }
    return binarySearchL(target, array, l, mid - 1, w)
}

func binarySearchN(target uint, array [NODE_NUM]uint, l, r, w int) (int, bool) {
    if l > r {
        if l >= w {
            return l, false
        }
        return r, false
    }
    mid := (l + r) >> 1
    if array[mid] == target {
        return mid, true
    }
    if array[mid] < target {
        return binarySearchN(target, array, mid + 1, r, w)
    }
    return binarySearchN(target, array, l, mid - 1, w)
}

func searchCSBTNode(node interface{}, key uint) (uint) {
    if l, ok := node.(*CSBTLeaf);ok {
        if i, b := binarySearchL(key, l.key, 0, l.keyNum - 1, l.keyNum);b {
            return l.data[i]
        }
        return 0
    }
    n := node.(*CSBTNode)
    i ,_ := binarySearchN(key, n.key, 0, n.keyNum - 1, n.keyNum)
    if cn, ok := n.child.(*CSBTLeaf);ok {
        n := (*CSBTLeaf)(unsafe.Pointer(uintptr(unsafe.Pointer(cn)) + unsafe.Sizeof(cn) * uintptr(i)))
        return searchCSBTNode(n, key)
    }
    cl := n.child.(*CSBTNode)
    l := (*CSBTNode)(unsafe.Pointer(uintptr(unsafe.Pointer(cl)) + unsafe.Sizeof(cl) * uintptr(i)))
    return searchCSBTNode(l, key)
}

//Return leafNode and key and ifFound
func selectCSBTNode(node interface{}, key uint) (*CSBTLeaf, int, bool) {
    if l, ok := node.(*CSBTLeaf);ok {
        i, b := binarySearchL(key, l.key, 0, l.keyNum - 1, l.keyNum)
        return l, i, b
    }
    n := node.(*CSBTNode)
    i ,_ := binarySearchN(key, n.key, 0, n.keyNum - 1, n.keyNum)
    if cn, ok := n.child.(*CSBTLeaf);ok {
        n := (*CSBTLeaf)(unsafe.Pointer(uintptr(unsafe.Pointer(cn)) + unsafe.Sizeof(cn) * uintptr(i)))
        return selectCSBTNode(n, key)
    }
    cl := n.child.(*CSBTNode)
    l := (*CSBTNode)(unsafe.Pointer(uintptr(unsafe.Pointer(cl)) + unsafe.Sizeof(cl) * uintptr(i)))
    return selectCSBTNode(l, key)
}

func bytes2node(b []byte) *CSBTNode {
    return (*CSBTNode)(unsafe.Pointer(
        (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
    ))
}

func bytes2leaf(b []byte) *CSBTLeaf {
    return (*CSBTLeaf)(unsafe.Pointer(
        (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
    ))
}

func insertToDCSBT(tree *DCSBT, key uint, value uint, base uintptr) (error) {
    l, k, b := selectCSBTNode(tree.indexHeader.root, key)
    if b {
        return errors.New("already exist key:" + string(key))
    }
    if l.keyNum < NODE_NUM / 2 {
        for i := l.keyNum;i > k;i-- {
            l.key[i] = l.key[i - 1]
        }
        l.key[k] = key
        l.data[k] = value
        l.base[k] = base
        l.keyNum++
        return nil
    }
    if root, ok := tree.indexHeader.root.(*CSBTLeaf);ok { //Root is leaf, key must in root
        sp1 := makeBytes(tree.freePos, (NODE_NUM / 3) * 24 + 24)
        splited1 := (*CSBTLeaf)(unsafe.Pointer(&sp1[0]))
        tree.freePos += unsafe.Sizeof(splited1)
        if unsafe.Sizeof(splited1) != (NODE_NUM / 3) * 24 + 24 {
            panic("serialize RAM address error!")
        }
        sp2 := makeBytes(tree.freePos, (NODE_NUM / 3) * 24 + 24)
        splited2 := (*CSBTLeaf)(unsafe.Pointer(&sp2[0]))
        tree.freePos += unsafe.Sizeof(splited2)
        splited1.keyNum = NODE_NUM / 6 
        splited2.keyNum = NODE_NUM / 3 - NODE_NUM / 6
        for i := 0;i < NODE_NUM / 6;i++ {
            splited1.key[i] = root.key[i]
            splited1.base[i] = root.base[i]
            splited1.data[i] = root.data[i]
        }
        for i := root.keyNum - 1;i > NODE_NUM/6 - 1;i-- {
            splited1.key[i] = root.key[i]
        }
        r := makeBytes(tree.freePos, NODE_NUM * 8 + 24)
        newRoot := (*CSBTNode)(unsafe.Pointer(&r))
        newRoot.keyNum = 2
        newRoot.key[0] = splited1.key[splited1.keyNum - 1]
        newRoot.key[1] = splited2.key[splited2.keyNum - 1]
        newRoot.child = splited1
        tree.indexHeader.root = newRoot
        splited1.key[splited1.keyNum] = key
        splited1.data[splited1.keyNum] = value
        splited1.base[splited1.keyNum] = base
        return nil
    }
    return nil
}