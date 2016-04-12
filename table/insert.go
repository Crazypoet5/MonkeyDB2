package table

func (t *Table) Insert(columnNames []string, data [][][]byte) {
	fields := t.Fields
	for _, row := range data {
		t.FirstPage.Append(uint2bytes(0))
		columnNamesP := 0
		for i := 0; i < len(fields); i++ {
			if columnNames == nil || len(columnNames) == 0 ||
				columnNamesP < len(columnNames) && columnNames[columnNamesP] == fields[i].Name {
				t.FirstPage.AppendField(&fields[i], row[columnNamesP])
				columnNamesP++
				if columnNamesP >= len(columnNames) {
					break
				}
			} else {
				t.FirstPage.AppendField(&fields[i], nil)
			}
		}

	}
}
