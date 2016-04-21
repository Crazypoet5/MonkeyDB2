package table

import (
	"unsafe"

	"../index"
)

//Dunplicated
func (t *Table) Insert_dunplicated(columnNames []string, data [][][]byte) {
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

func (t *Table) Insert(fieldMap map[int]int, data [][][]byte) {
	fields := t.Fields
	for _, row := range data {
		pos := t.LastPage.GetFreePos()
		t.LastPage.Append(uint2bytes(0)) //Skip
		for i := 0; i < len(fields); i++ {
			if k, ok := fieldMap[i]; ok {
				t.LastPage.AppendField(&fields[i], row[k])
				if fields[i].Index != nil {
					var p uint
					p = uint(uintptr(unsafe.Pointer(t.LastPage))) << 24
					p = p | pos
					fields[i].Index.Insert(index.BKDRHash(row[k]), uintptr(p))
				}
			} else {
				t.LastPage.AppendField(&fields[i], nil)
			}
		}

	}
}
