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
		currentPtr:  32,
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

func (p *Reader) DumpPage() *exe.Relation {
	oldP := p.currentPtr
	p.currentPtr = 32
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
