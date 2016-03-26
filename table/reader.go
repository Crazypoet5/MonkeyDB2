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
    ret := make(Relation)
    table := p.currentPage.GetTable()
    row := make([]Value)
    for f := range table.Fields {
        if f.FixedSize {
            v := p.currentPage.Read(p.currentPtr, f.Size)
            ret = append(row, v)
            p.currentPtr += f.Size
        } else {
            size := bytes2uint32(p.currentPage.Read(p.currentPtr, 4))
            p.currentPtr += 4
            v := p.currentPage.Read(p.currentPtr, int(size))
            ret = append(row, v)
             p.currentPtr += int(size)
        }
    }
    ret = append(ret, row)
    return ret
}