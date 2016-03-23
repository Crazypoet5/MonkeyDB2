package mem

import (
    "../../memory"
)

func LoadManagedBlockFromDataBlock(db *memory.DataBlock) *ManagedBlock {
    return &ManagedBlock {
        db:         db,
    }
}

func LoadManagedBlockFromOldUintptr(p uintptr) *ManagedBlock {
    newPtr := uintptr(0)
    for k, v := range memory.RecoveryTable {
        if k == p {
            newPtr = v
            break
        }
    }
    if newPtr == 0 {
        return nil
    }
    for l := memory.DataBlockList.Front();l != nil;l = l.Next() {
        db := l.Value.(*memory.DataBlock)
        if db.RawPtr == newPtr {
            return &ManagedBlock {
                db:         db,
            }
        }
    }
    return nil
}