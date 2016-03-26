package table

type Reader struct {
    currentPtr      uint
    currentPage     *Page
}

func (p *Page) NewReader() *Reader {
    return &Reader {
        currentPtr:     24,
        currentPage:    p,
    }
}

func (p *Reader) PeekRecord() Relation {
    ret := make(Relation, 0)
    table := p.currentPage.GetTable()
    row := make([]Value, 0)
    for _, f := range table.Fields {
        if f.FixedSize {
            v, _ := p.currentPage.Read(p.currentPtr, uint(f.Size))
            row = append(row, v)
            p.currentPtr += uint(f.Size)
        } else {
            data, _ := p.currentPage.Read(p.currentPtr, 4)
            size := bytes2uint32(data)
            p.currentPtr += 4
            v, _ := p.currentPage.Read(p.currentPtr, uint(size))
            row = append(row, v)
            p.currentPtr += uint(size)
        }
    }
    ret = append(ret, row)
    return ret
}