package mem

func (mb *ManagedBlock) Copy(dst uint, src uint, length int) {
	data, _ := mb.Read(src, uint(length))
	mb.Write(dst, data)
}
