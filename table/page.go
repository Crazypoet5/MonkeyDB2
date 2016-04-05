package table

import (
	"unsafe"

	"../exe"
	"../memory"
)

const (
	NORMAL_PAGE_SIZE = 1024
	PREV_OFFSET      = 8
	NEXT_OFFSET      = 16
	FREE_P_OFFSET    = 24
)

type Page struct {
	memory.DataBlock
}

// This struct is only to refrence to programer
type page struct {
	table uintptr
	prev  uintptr
	next  uintptr
	fp    uint
}

func (p *Page) GetTable() *Table {
	ptr, _ := p.Read(0, 8)
	return (*Table)(unsafe.Pointer(uintptr(bytes2uint(ptr))))
}

func NewPage() *Page {
	db, _ := memory.CreateImage(NORMAL_PAGE_SIZE)
	db.Write(PREV_OFFSET, uint2bytes(0))
	db.Write(NEXT_OFFSET, uint2bytes(0))
	db.Write(FREE_P_OFFSET, uint2bytes(32))
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

func (p *Page) GetFreePos() uint {
	data, _ := p.Read(FREE_P_OFFSET, 8)
	return bytes2uint(data)
}

func (p *Page) ForwardFreePos(i uint) {
	fp := p.GetFreePos()
	fp += i
	p.Write(FREE_P_OFFSET, uint2bytes(fp))
}

func (p *Page) AppendField(f *Field, data []byte) {
	if data == nil || len(data) == 0 {
		switch f.Type {
		case exe.INT, exe.FLOAT, exe.OBJECT, exe.ARRAY:
			data = make([]byte, 8)
		case exe.STRING:
			data = make([]byte, 0)
		}
	}
	if !f.FixedSize {
		p.Append(uint2bytes(uint(len(data))))
	}
	p.Append(data)
}

func (p *Page) Append(data []byte) {
	fp := p.GetFreePos()
	p.Write(fp, data)
	fp += uint(len(data))
	p.Write(FREE_P_OFFSET, uint2bytes(fp))
}
