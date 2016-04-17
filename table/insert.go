package table

func (t *Table) Insert(columnNames []string, data [][][]byte) {
	fields := t.Fields
	for _, row := range data {
		t.LastPage.Append(uint2bytes(0)) //Skip
		columnNamesP := 0
		for i := 0; i < len(fields); i++ {
			if columnNames == nil || len(columnNames) == 0 ||
				columnNamesP < len(columnNames) && columnNames[columnNamesP] == fields[i].Name {
				t.LastPage.AppendField(&fields[i], row[columnNamesP])
				columnNamesP++
				if len(columnNames) != 0 && columnNames != nil && columnNamesP >= len(columnNames) {
					break
				}
			} else {
				t.LastPage.AppendField(&fields[i], nil)
			}
		}

	}
}
