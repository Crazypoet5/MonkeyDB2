package main

import (
    "fmt"
    "unsafe"
)

type node struct {
    isLeaf      byte
    reversed    byte
    keyNum      uint16
    key         [13]uint32
    child       uintptr
}

type leaf struct {
    isLeaf      byte
    keyNum      byte
    key         [3]uint32
    value       [3]uintptr
    left,right  uintptr
    reverse     uint32
    reverse2    uint16
}

func main() {
    var n node
    var l leaf
    fmt.Println(unsafe.Sizeof(n), unsafe.Sizeof(l))
}