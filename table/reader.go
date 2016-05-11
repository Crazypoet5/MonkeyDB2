package table

import (
	"../exe"
	"../index"
)

var RECORD_MAX_LENGTH = 4096

type Reader struct {
	currentPtr  uint
	currentPage *Page
}

func (p *Page) NewReader() *Reader {
	return &Reader{
		currentPtr:  64,
		currentPage: p,
	}
}

func (p *Reader) PeekRecord() *exe.Relation {
	//	if p.currentPtr == p.currentPage.GetEOP() {
	//		p.currentPage = p.currentPage.NextPage()
	//		p.currentPtr = 64
	//	}
	v, _ := p.currentPage.Read(p.currentPtr, 8)
	p.currentPtr += 8
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

func (p *Reader) NextRecord() {
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
			p.currentPtr += uint(f.Size)
		} else {
			data, _ := p.currentPage.Read(p.currentPtr, 4)
			size := bytes2uint32(data)
			p.currentPtr += 4
			p.currentPtr += uint(size)
		}
	}
}

func (p *Reader) DumpPage() *exe.Relation {
	oldP := p.currentPtr
	p.currentPtr = 64
	ret := exe.NewRelation()
	tab := p.currentPage.GetTable()
	columns := make([]string, 0)
	for i := 0; i < len(tab.Fields); i++ {
		columns = append(columns, tab.Fields[i].Name)
	}
	ret.SetColumnNames(columns)

	//	v, _ := p.currentPage.Read(0, 512)
	//	fmt.Println(v)
	for p.currentPtr < p.currentPage.GetFreePos() {
		v, _ := p.currentPage.Read(p.currentPtr, 8)
		p.currentPtr += 8
		skip := bytes2uint(v)
		//		fmt.Println(skip)
		if skip != 0 {
			p.currentPtr += skip
			continue
		}
		table := p.currentPage.GetTable()
		row := make(exe.Row, 0)
		for _, f := range table.Fields {
			if f.FixedSize {
				v, _ := p.currentPage.Read(p.currentPtr, uint(f.Size))
				//				fmt.Println(v)
				row = append(row, exe.NewValue(f.Type, v))
				p.currentPtr += uint(f.Size)
			} else {
				data, _ := p.currentPage.Read(p.currentPtr, 4)
				size := bytes2uint32(data)
				//				fmt.Println(size)
				p.currentPtr += 4
				v, _ := p.currentPage.Read(p.currentPtr, uint(size))
				//				fmt.Println(v)
				row = append(row, exe.NewValue(f.Type, v))
				p.currentPtr += uint(size)
			}
		}
		ret.AddRow(row)
	}
	p.currentPtr = oldP
	return ret
}

func (p *Reader) DumpTable() *exe.Relation {
	oldPage := p.currentPage
	oldPtr := p.currentPtr
	ret := exe.NewRelation()
	for p.currentPage != nil {
		p.currentPtr = 0
		r := p.DumpPage()
		ret.Rows = append([]exe.Row{}, r.Rows...)
		ret.SetColumnNames(r.ColumnNames)
		p.currentPage = p.currentPage.NextPage()
	}
	p.currentPage = oldPage
	p.currentPtr = oldPtr
	return ret
}

func (p *Reader) LoadIndex(i *index.Index, kField int) error {
	oldPage := p.currentPage
	oldPtr := p.currentPtr
	for p.currentPage != nil {
		p.currentPtr = 64
		for p.currentPtr < p.currentPage.GetFreePos() {
			pos := p.currentPage.RawPtr
			v, _ := p.currentPage.Read(p.currentPtr, 8)
			p.currentPtr += 8
			skip := bytes2uint(v)
			if skip != 0 {
				p.currentPtr += skip
				continue
			}
			ptr := uint(pos)
			ptr <<= 24
			ptr |= p.currentPtr
			table := p.currentPage.GetTable()
			for k, f := range table.Fields {
				if f.FixedSize {
					if kField == k {
						v, _ := p.currentPage.Read(p.currentPtr, uint(f.Size))
						err := i.I.Insert(index.BKDRHash(v), uintptr(ptr))
						if err != nil {
							return err
						}
					}
					p.currentPtr += uint(f.Size)
				} else {
					data, _ := p.currentPage.Read(p.currentPtr, 4)
					size := bytes2uint32(data)
					p.currentPtr += 4
					if kField == k {
						v, _ := p.currentPage.Read(p.currentPtr, uint(size))
						err := i.I.Insert(index.BKDRHash(v), uintptr(ptr))
						if err != nil {
							return err
						}
					}
					p.currentPtr += uint(size)
				}
			}
		}
		p.currentPage = p.currentPage.NextPage()
	}
	p.currentPage = oldPage
	p.currentPtr = oldPtr
	return nil
}
