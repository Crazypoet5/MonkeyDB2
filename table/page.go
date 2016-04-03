package table

import (
	"unsafe"

	"../memory"
)

const (
	NORMAL_PAGE_SIZE = 1024
	PREV_OFFSET      = 8
	NEXT_OFFSET      = 16
)

type Page struct {
	memory.DataBlock
}

// This struct is only to refrence to programer
type page struct {
	table uintptr
	prev  uintptr
	next  uintptr
}

func (p *Page) GetTable() *Table {
	ptr, _ := p.Read(0, 8)
	return (*Table)(unsafe.Pointer(uintptr(bytes2uint(ptr))))
}

func NewPage() *Page {
	db, _ := memory.CreateImage(NORMAL_PAGE_SIZE)
	db.Write(PREV_OFFSET, uint2bytes(0))
	db.Write(NEXT_OFFSET, uint2bytes(0))
	return &Page{
		DataBlock: *db,
	}
}

func (p *Page) NextPage() *Page {
	ptr, _ := p.Read(NEXT_OFFSET, 8)
	return (*Page)(unsafe.Pointer(uintptr(bytes2uint(ptr))))
}

func (p *Page) PrevPage() *Page {
	ptr, _ := p.Read(PREV_OFFSET, 8)
	return (*Page)(unsafe.Pointer(uintptr(bytes2uint(ptr))))
}

func (p *Page) AppendPage() {
	pNew := NewPage()
	p.Write(NEXT_OFFSET, uint2bytes(uint(uintptr(unsafe.Pointer(pNew)))))
	pNew.Write(PREV_OFFSET, uint2bytes(uint(uintptr(unsafe.Pointer(p)))))
}
