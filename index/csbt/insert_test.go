package csbt

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	tree := NewDCSBT()
	for i := 0; i < 100; i++ {
		tree.Insert(uint32(i), uintptr(i))
	}
	nd95 := tree.FindLeaf(94, false)
	l95 := tree.MB.GetChild(nd95.node, nd95.offset)
	l95Next := tree.MB.GetLeafRight(l95)
	fmt.Println(tree.MB.GetLeafKey(l95Next, 0))
}
