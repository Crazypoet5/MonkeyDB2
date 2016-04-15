package mem

// nSize start by 1
func (mb *ManagedBlock) Malloc(nSize int) uint {
	ret := mb.GetFreePos()
	mb.AddFreePos(uint(nSize) * 64)
	return ret
	// TODO:    add copy and realloc
	return 0
}
