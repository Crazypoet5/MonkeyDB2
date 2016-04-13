package mem

import (
	"../../memory"
)

func LoadManagedBlockFromDataBlock(db *memory.DataBlock) *ManagedBlock {
	return &ManagedBlock{
		DataBlock: *db,
	}
}

func LoadManagedBlockFromOldUintptr(p uintptr) *ManagedBlock {
	return &ManagedBlock{
		DataBlock: *memory.RecoveryTable[p],
	}
}
