package table

import (
	"unsafe"

	"../memory"
)

// Recovery page list and return the last page
func (p *Page) Recovery(t *Table, prev *Page) *Page {
	p.Write(0, uint2bytes(uint(uintptr(unsafe.Pointer(t)))))
	p.Write(PREV_OFFSET, uint2bytes(uint(uintptr(unsafe.Pointer(prev)))))
	ptr, _ := p.Read(NEXT_OFFSET, 8)
	db := memory.RecoveryTable[uintptr(bytes2uint(ptr))]
	if db != nil {
		pNext := &Page{
			DataBlock: *db,
		}
		p.Write(NEXT_OFFSET, uint2bytes(uint(uintptr(unsafe.Pointer(pNext)))))
		return pNext.Recovery(t, p)
	}
	return p
}
