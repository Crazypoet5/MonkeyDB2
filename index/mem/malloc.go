package mem

// nSize start by 1
func (mb *ManagedBlock) Malloc(nSize int) uint {
    for i := nSize;i <= MAX_POWER;i++ {
        ret := mb.ListPopFront(i)
        if ret != 0 {
            return ret
        }
    }
    // TODO:    add copy and realloc
    return 0
}