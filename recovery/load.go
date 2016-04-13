package recovery

import (
	"../memory"
	"../table"
)

func LoadData() {
	memory.Recovery()
	table.Recovery()
}
