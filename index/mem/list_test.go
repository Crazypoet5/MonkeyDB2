package mem

import (
    "testing"
    "fmt"
)

func TestNewManagedBlock(t *testing.T) {
    mb := NewManagedBlock()
    for i := 0;i < MAX_POWER;i++ {
        fmt.Println(mb.freeList[i].Len())
    }
    // memory.SyncAllImageToFile()
}