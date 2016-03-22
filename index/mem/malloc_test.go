package mem

import (
    "testing"
)

func TestMalloc(t *testing.T) {
    mb := NewManagedBlock()
    mb.Malloc(12)
}