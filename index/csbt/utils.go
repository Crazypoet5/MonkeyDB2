package csbt

import (
    "unsafe"
    "reflect"
)

func makeBytes(p uintptr, n int) []byte {
    var b []byte
    header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    header.Cap = n
    header.Len = n
    header.Data = p
    return b
}

func (tree *DCSBT) allocate(n int) []byte {
    b := makeBytes(tree.freePos, n)
    tree.freePos += uintptr(n)
    return b
}

func (parent *CSBTNode) children(i int) interface{} {
    p := unsafe.Pointer(uintptr(unsafe.Pointer(parent.child.(*CSBTLeaf))) + uintptr(i))
    if _, ok := parent.child.(*CSBTLeaf);ok {
        return (*CSBTLeaf)(p)
    }
    return (*CSBTNode)(p)
}