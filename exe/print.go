package exe

import (
	"fmt"
)

func (r *Relation) Print() {
	if r == nil || len(r.Rows) == 0 {
		return
	} else {
		fmt.Println(r.Rows)
	}
}
