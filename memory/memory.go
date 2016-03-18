package memory

import (
    "sync"
)

var ImageTable = make(map[uintptr]string)

type DataBlock struct {
    RawPtr      uintptr
    Size        int
    RWMutex     sync.RWMutex
}
