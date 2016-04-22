package table

func (t *Table) Drop() {
	for _, v := range t.Fields {
		if v.Index != nil {
			v.Index.Delete()
		}
	}
	for p := t.FirstPage; p != nil; {
		pNext := p.NextPage()
		p.Delete()
		p = pNext
	}
	for k, v := range TableList {
		if v == t {
			TableList = append(TableList[0:k], TableList[k+1:]...)
			return
		}
	}
}
