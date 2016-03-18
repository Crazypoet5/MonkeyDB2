package memory

import (
    "sync"
)

var ImageTable = make(map[uintptr]string)
var RecoveryTable = make(map[uintptr]uintptr)

type DataBlock struct {
    RawPtr      uintptr
    Size        int
    RWMutex     sync.RWMutex
}
