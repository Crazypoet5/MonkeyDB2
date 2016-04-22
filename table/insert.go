package table

import (
	"strconv"
	"time"

	"../common"

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

//Map Table -> Data
func (t *Table) Insert(fieldMap map[int]int, data [][][]byte) error {
	fields := t.Fields
	for _, row := range data {
		pos := t.LastPage.GetFreePos()
		var indexLog = make(map[index.Indexer]uint32)
		t.LastPage.Append(uint2bytes(0)) //Skip
		for i := 0; i < len(fields); i++ {
			if k, ok := fieldMap[i]; ok {
				if fields[i].Index != nil {
					var p uint
					p = uint(t.LastPage.RawPtr) << 24
					p = p | pos
					start := time.Now().UnixNano()
					key := index.BKDRHash(row[k])
					err := fields[i].Index.I.Insert(key, uintptr(p))
					if err != nil {
						for k, v := range indexLog {
							k.Delete(v)
						}
						t.LastPage.SetFreePos(pos)
						return err
					}
					common.Print2("Insert in Index within " + strconv.Itoa(int(time.Now().UnixNano()-start)) + " ns.")
					indexLog[fields[i].Index.I] = key
				}
				t.LastPage.AppendField(&fields[i], row[k])
			} else {
				t.LastPage.AppendField(&fields[i], nil)
			}
		}

	}
	return nil
}
