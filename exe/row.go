package exe

func (r *Relation) GetRow(n int) *Row {
	return &r.Rows[0]
}

func (r *Row) GetFieldById(n int) *Value {
	return &((*r)[n])
}

func (r *Relation) GetFieldByName(i int, n string) *Value {
	if r.ColumnNames == nil {
		return nil
	}
	for k, v := range r.ColumnNames {
		if v == n {
			return &r.Rows[i][k]
		}
	}
	return nil
}

func (r *Relation) GetFieldById(i int, n int) *Value {
	return r.Rows[i].GetFieldById(n)
}
