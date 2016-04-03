package table

import (
	"../exe"
)

type Reader struct {
	currentPtr  uint
	currentPage *Page
}

func (p *Page) NewReader() *Reader {
	return &Reader{
		currentPtr:  24,
		currentPage: p,
	}
}

func (p *Reader) PeekRecord() *exe.Relation {
	v, _ := p.currentPage.Read(p.currentPtr, 8)
	skip := bytes2uint(v)
	if skip != 0 {
		p.currentPtr += skip
		return p.PeekRecord()
	}
	ret := exe.NewRelation()
	table := p.currentPage.GetTable()
	row := make(exe.Row, 0)
	for _, f := range table.Fields {
		if f.FixedSize {
			v, _ := p.currentPage.Read(p.currentPtr, uint(f.Size))
			row = append(row, exe.NewValue(f.Type, v))
			p.currentPtr += uint(f.Size)
		} else {
			data, _ := p.currentPage.Read(p.currentPtr, 4)
			size := bytes2uint32(data)
			p.currentPtr += 4
			v, _ := p.currentPage.Read(p.currentPtr, uint(size))
			row = append(row, exe.NewValue(f.Type, v))
			p.currentPtr += uint(size)
		}
	}
	ret.AddRow(row)
	return ret
}
