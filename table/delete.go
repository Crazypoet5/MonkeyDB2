package table

import (
	"log"

	"../exe"
	"../index"
)

func (t *Table) Delete(ids *exe.BitSet) {
	reader := t.FirstPage.NewReader()
	for i := 0; i < ids.Len(); i++ {
		pos := reader.currentPtr
		page := reader.currentPage

		if ids.Get(i) {
			reader.DeleteRecordIndex()
			var skip uint
			if page == reader.currentPage {
				skip = reader.currentPtr - pos - 8
				page.Write(pos, uint2bytes(skip))
				log.Println("in-page delete")
			} else {
				skip = reader.currentPtr - 72
				reader.currentPage.Write(64, uint2bytes(skip))
				log.Println("out-page delete")
			}

		} else {
			reader.NextRecord()
		}
	}
}

func (p *Page) DeleteFromOffset(offset uint) {
	reader := p.NewReader()
	reader.currentPtr = offset
	reader.NextRecord()
	end := reader.currentPtr
	skip := end - offset
	p.Write(offset, uint2bytes(skip))
}

//Delete current record Index and goto next record
func (p *Reader) DeleteRecordIndex() {
	p.CheckPage()
	if p.currentPage == nil {
		log.Println("attempt to delete from nil page")
		return
	}
	v, _ := p.currentPage.Read(p.currentPtr, 8)
	p.currentPtr += 8
	skip := bytes2uint(v)
	if skip != 0 {
		p.currentPtr += skip
		p.NextRecord()
		return
	}
	table := p.currentPage.GetTable()
	for _, f := range table.Fields {
		if f.FixedSize {
			v, _ := p.currentPage.Read(p.currentPtr, uint(f.Size))
			if f.Index != nil {
				f.Index.I.Delete(index.BKDRHash(v))
			}
			p.currentPtr += uint(f.Size)
		} else {
			data, _ := p.currentPage.Read(p.currentPtr, 4)
			size := bytes2uint32(data)
			p.currentPtr += 4
			v, _ := p.currentPage.Read(p.currentPtr, uint(size))
			if f.Index != nil {
				f.Index.I.Delete(index.BKDRHash(v))
			}
			p.currentPtr += uint(size)
		}
	}
}
