package recovery

import (
	"../index"
	"../memory"
	"../table"
)

func LoadData() {
	memory.Recovery()
	index.Recovery()
	table.Recovery()
}
