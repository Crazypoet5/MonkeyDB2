package table

import (
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
			skip := reader.currentPtr - pos - 8
			page.Write(pos, uint2bytes(skip))
		} else {
			reader.NextRecord()
		}
	}
}

//Delete current record Index and goto next record
func (p *Reader) DeleteRecordIndex() {
	//	if p.currentPtr == p.currentPage.GetEOP() {
	//		p.currentPage = p.currentPage.NextPage()
	//		p.currentPtr = 64
	//		p.NextRecord()
	//		return
	//	}
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
