//      DCSBT Index @InsZVA 2015
//      Based on WANG Sheng, QIN Xiaolin, SHEN Yao, et al. Research on durable CSB+-tree indexing technology. Journal of Frontiers of Computer Science and Technology, 2015, 9(2): 182-192.

package csbt

import (
	"../mem"
)

type leaf struct {
	isLeaf      byte
	keyNum      byte
	reserved    uint16
	key         [3]uint32
	value       [3]uintptr
	left, right uint
	reserved2   uint
}

type Node struct {
	IsLeaf   byte
	Reversed byte
	KeyNum   uint16
	Key      [13]uint32
	Child    uint
}

type Header struct {
	Root     uint
	Min, Max uint
}

type DCSBT struct {
	MB *mem.ManagedBlock
}

func NewDCSBT() *DCSBT {
	t := &DCSBT{
		MB: mem.NewManagedBlock(),
	}
	leaf := t.MB.NewLeaves(1)
	t.MB.SetRoot(leaf)
	t.MB.SetMin(leaf)
	t.MB.SetMax(leaf)
	return t
}
