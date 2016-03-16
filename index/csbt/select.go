package csbt

import (
    "unsafe"
)

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

//Return leafNode and key and its parent node and ifFound
func selectCSBTNode(node interface{}, key uint, parent interface{}) (*CSBTLeaf, int, bool, *CSBTNode) {
    if l, ok := node.(*CSBTLeaf);ok {
        i, b := binarySearchL(key, l.key, 0, l.keyNum - 1, l.keyNum)
        if p, ok := parent.(*CSBTNode);ok {
            return l, i, b, p
        }
        return l, i, b, nil
    }
    n := node.(*CSBTNode)
    i ,_ := binarySearchN(key, n.key, 0, n.keyNum - 1, n.keyNum)
    if cn, ok := n.child.(*CSBTLeaf);ok {
        n := (*CSBTLeaf)(unsafe.Pointer(uintptr(unsafe.Pointer(cn)) + unsafe.Sizeof(cn) * uintptr(i)))
        return selectCSBTNode(n, key, n)
    }
    cl := n.child.(*CSBTNode)
    l := (*CSBTNode)(unsafe.Pointer(uintptr(unsafe.Pointer(cl)) + unsafe.Sizeof(cl) * uintptr(i)))
    return selectCSBTNode(l, key, n)
}
