package index

func BKDRHash(bytes []byte) uint32 {
	seed := uint32(13131313)
	hash := uint32(0)
	for i := 0; i < len(bytes); i++ {
		hash = hash*seed + uint32(bytes[i])
	}
	return hash & 0x7fffffff
}
