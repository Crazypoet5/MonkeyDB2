package mem

import (
    "testing"
    "fmt"
)

var mb = NewManagedBlock()

func TestNewManagedBlock(t *testing.T) {
    for i := 0;i < MAX_POWER;i++ {
        fmt.Println(mb.ListGetLength(i))
    }
    // memory.SyncAllImageToFile()
}

func TestListPopBack(t *testing.T) {
    old := mb.ListGetLength(1)
    mb.ListPopBack(1)
    if old - mb.ListGetLength(1) != 1 {
        t.Error(old)
    }
    for i := 0;i < old - 2;i++ {
        (mb.ListPopBack(1))
    }
    if mb.ListGetBack(mb.GetFreeList(1)) == 0 {
        t.Error(mb.ListGetBack(mb.GetFreeList(1)))
    }
    mb.ListPopBack(1)
    if mb.ListGetBack(mb.GetFreeList(1)) != 0 {
        t.Error(mb.ListGetBack(mb.GetFreeList(1)))
    }
}

func TestListPopFront(t *testing.T) {
    mb = NewManagedBlock()
    old := mb.ListGetLength(1)
    mb.ListPopFront(1)
    if old - mb.ListGetLength(1) != 1 {
        t.Error(old)
    }
    for i := 0;i < old - 2;i++ {
        (mb.ListPopFront(1))
    }
    if mb.ListGetBack(mb.GetFreeList(1)) == 0 {
        t.Error(mb.ListGetBack(mb.GetFreeList(1)))
    }
    mb.ListPopFront(1)
    if mb.ListGetBack(mb.GetFreeList(1)) != 0 {
        t.Error(mb.ListGetBack(mb.GetFreeList(1)))
    }
}