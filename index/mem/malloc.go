package mem

// nSize start by 1
func (mb *ManagedBlock) Malloc(nSize int) uint {
    for i := uint(nSize);i <= MAX_POWER;i++ {
        if mb.freeList[i - 1].Len() != 0 {
            continue
        }
        p := mb.freeList[i - 1].Front()
        mb.freeList[i - 1].Remove(p)
        mb.used += 64 * i
        return p.Value.(uint)
    }
    // TODO:    add copy and realloc
    return 0
}