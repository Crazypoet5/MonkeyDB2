package csbt

import (
	"../cursor"
)

type Cursor struct {
	tree   *DCSBT
	leaf   uint
	offset int
}

func (c *Cursor) Next() cursor.Cursor {
	keyNum := c.tree.MB.GetLeafKeyNum(c.leaf)
	if c.offset < keyNum {
		return &Cursor{
			tree:   c.tree,
			leaf:   c.leaf,
			offset: c.offset + 1,
		}
	}
	l := c.tree.MB.GetLeafRight(c.leaf)
	if l == 0 {
		return nil
	}
	return &Cursor{
		tree:   c.tree,
		leaf:   l,
		offset: 0,
	}
}

func (c *Cursor) Prev() cursor.Cursor {
	if c.offset > 0 {
		return &Cursor{
			tree:   c.tree,
			leaf:   c.leaf,
			offset: c.offset - 1,
		}
	}
	l := c.tree.MB.GetLeafLeft(c.leaf)
	if l == 0 {
		return nil
	}
	return &Cursor{
		tree:   c.tree,
		leaf:   l,
		offset: 0,
	}
}

// Return uintptr or 0 if deleted
func (c *Cursor) Read() (uintptr, uint) {
	i := uint(c.tree.MB.GetLeafValue(c.leaf, c.offset))
	p := uint(uintptr(i >> 24))
	offset := i & 0x0000000000111111 // < 1MB
	return uintptr(p), offset
}

func (c *Cursor) Delete() {
	c.tree.MB.SetLeafValue(c.leaf, c.offset, 0)
}

func (c *Cursor) Write(p uintptr, offset uint) {
	v := uint(0)
	v = uint(p) << 24
	v = v | offset
	c.tree.MB.SetLeafValue(c.leaf, c.offset, uintptr(v))
}
