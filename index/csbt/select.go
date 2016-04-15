package csbt

import (
	"../cursor"
)

// Select will return a cursor if found and return nil else
func (t *DCSBT) Select(k uint32) cursor.Cursor {
	var l uint
	n := t.FindLeaf(k, false)
	if n.node == 0 {
		l = t.MB.GetRoot()
	} else {
		l = t.MB.GetChild(n.node, n.offset)
	}
	if l == 0 {
		return nil
	}
	for i := 0; i < t.MB.GetLeafKeyNum(l); i++ {
		if t.MB.GetLeafKey(l, i) == k {
			return &Cursor{
				tree:   t,
				leaf:   l,
				offset: i,
			}
		}
	}
	return nil
}
