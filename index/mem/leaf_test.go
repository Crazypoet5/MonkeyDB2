package mem

import (
    "testing"
)

func TestNewLeaves(t *testing.T) {
    mb := NewManagedBlock()
    leaf := mb.NewLeaves(1)
    mb.SetLeafKeyNum(leaf, 1)
    mb.SetLeafKey(leaf, 0, 156)
    if mb.GetLeafKey(leaf, 0) != 156 {
        t.Error(mb.GetLeafKey(leaf, 0))
    }
}