package exe

func (r *Relation) Project(columnNames []string) *Relation {

	ids := make([]int, 0)
	for i := 0; i < len(columnNames); i++ {
		for k, v := range r.ColumnNames {
			if v == columnNames[i] {
				ids = append(ids, k)
			}
		}
	}
	newR := NewRelation()
	newR.SetColumnNames(columnNames)
	for _, r := range r.Rows {
		row := []Value{}
		for i := 0; i < len(ids); i++ {
			row = append(row, r[ids[i]])
		}
		newR.AddRow(row)
	}
	return newR
}
