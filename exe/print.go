package exe

import (
	"fmt"
	"strconv"
	"unsafe"
)

func (v *Value) Print() {
	switch v.Kind {
	case INT:
		i := *(*int)(unsafe.Pointer(&v.Raw[0]))
		fmt.Print(i)
	case FLOAT:
		f := *(*float64)(unsafe.Pointer(&v.Raw[0]))
		fmt.Print(f)
	case STRING:
		s := (string)(v.Raw)
		fmt.Print(s)
	}
}

func (v *Value) StrLen() int {
	switch v.Kind {
	case INT:
		i := *(*int)(unsafe.Pointer(&v.Raw[0]))
		return len([]byte(strconv.Itoa(i)))
	case FLOAT:
		f := *(*float64)(unsafe.Pointer(&v.Raw[0]))
		return len([]byte(strconv.FormatFloat(f, 'f', -1, 64)))
	case STRING:
		return len(v.Raw)
	}
	return 8
}

func (r *Relation) Print() {

	if r == nil {
		return
	}
	var maxLen []int
	var colN int

	if r.ColumnNames != nil && len(r.ColumnNames) != 0 {
		colN = len(r.ColumnNames)
	} else if r.Rows != nil && len(r.Rows) != 0 {
		colN = len(r.Rows[0])
	} else {
		return
	}
	if r.Rows != nil && len(r.Rows) != 0 {
		maxLen = make([]int, colN)
		for i := 0; i < len(r.Rows); i++ {
			for j := 0; j < colN; j++ {
				if l := r.Rows[i][j].StrLen(); l > maxLen[j] {
					maxLen[j] = l
				}
			}
		}

	}
	if r.ColumnNames != nil && len(r.ColumnNames) != 0 {
		if maxLen == nil || len(maxLen) == 0 {
			maxLen = make([]int, colN)
		}
		for i := 0; i < colN; i++ {
			if l := len([]byte(r.ColumnNames[i])); l > maxLen[i] {
				maxLen[i] = l
			}
		}
	}
	if r.ColumnNames != nil && len(r.ColumnNames) != 0 {
		fmt.Print("_")
		for i := 0; i < colN; i++ {
			for j := 0; j < maxLen[i]; j++ {
				fmt.Print("_")
			}
			fmt.Print("_")
		}
		fmt.Print("\n")

		fmt.Print("|")
		for i := 0; i < colN; i++ {
			fmt.Printf("%s", r.ColumnNames[i])
			for k := 0; k < maxLen[i]-len([]byte(r.ColumnNames[i])); k++ {
				fmt.Print(" ")
			}
			fmt.Print("|")
		}
		fmt.Print("\n|")
		for i := 0; i < colN; i++ {
			for j := 0; j < maxLen[i]; j++ {
				fmt.Print("-")
			}
			fmt.Print("|")
		}
		fmt.Print("\n")
	}

	if r != nil && len(r.Rows) != 0 {
		for rw := 0; rw < len(r.Rows); rw++ {
			fmt.Print("|")
			for i := 0; i < colN; i++ {
				r.Rows[rw][i].Print()
				for k := 0; k < maxLen[i]-r.Rows[rw][i].StrLen(); k++ {
					fmt.Print(" ")
				}
				fmt.Print("|")
			}
			fmt.Print("\n|")
			for i := 0; i < colN; i++ {
				for j := 0; j < maxLen[i]; j++ {
					fmt.Print("-")
				}
				fmt.Print("|")
			}
			fmt.Print("\n")
		}
	}

	//	for i := 0; i < colN; i++ {
	//		for j := 0; j < maxLen[i]; j++ {
	//			fmt.Print("-")
	//		}
	//		fmt.Print("|")
	//	}
	//	fmt.Print("|\n")
}
