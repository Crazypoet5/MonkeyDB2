package index

func BKDRHash(bytes []byte) uint {
    seed := uint(13131313)
    hash := uint(0)
    for i := 0;i < len(bytes);i++ {
        hash = hash * seed + uint(bytes[i])
    } 
    return hash & 0x7fffffff
}