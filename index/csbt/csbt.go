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
    NODE_NUM = (CACHE_LINE - 16) / 8
    ALLOCATE = 2048
)

type DCSBT struct {
    root interface{}
    data []byte
    free int
}

type CSBTNode struct {
    keyNum int
    child interface{}
    key [NODE_NUM]uint
}

type CSBTLeaf struct {
    keyNum int
    key [NODE_NUM / 2]uint
    data [NODE_NUM / 2]uint
    left, right *CSBTLeaf
}

func NewCSBT() *CSBTLeaf {
    return &CSBTLeaf {}
}

func binarySearchL(target uint, array [NODE_NUM / 2]uint, l, r, w int) (int, bool) {
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

func insertToFree(tree *interface{}, key uint, value uint) (error) {
    l, k, b := selectCSBTNode(*tree, key)
    if b {
        return errors.New("already exist key:" + string(key))
    }
    if l.keyNum < NODE_NUM / 2 {
        for i := l.keyNum;i > k;i-- {
            l.key[i] = l.key[i - 1]
        }
        l.key[k] = key
        l.data[k] = value
        l.keyNum++
        return nil
    }
    if root, ok := (*tree).(*CSBTLeaf);ok { //Root is leaf, key must in root
        splited := 
        splited[0] := &CSBTLeaf {
            keyNum:     NODE_NUM / 4,
        }
        splited[1] := &CSBTLeaf {
            keyNum:     NODE_NUM / 2 - NODE_NUM / 4,
        }
        for i := 0;i < NODE_NUM/4;i++ {
            splited[0].key[i] = root.key[i]
        }
        for i := root.keyNum - 1;i > NODE_NUM/4 - 1;i-- {
            splited[1].key[i] = root.key[i]
        }
        
        newRoot := &CSBTNode {
            keyNum:     2,
        }
        newRoot.child
    }
}