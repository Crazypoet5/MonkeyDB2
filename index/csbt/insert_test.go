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
	fmt.Println(int(tree.Find(95)))
}
