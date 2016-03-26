package table

func uint162bytes(u uint16) []byte {
    b := make([]byte, 2)
    for i := 0;i < 2;i++ {
         b[i] = (byte)(u)
         u >>= 8
    }
    return b
}

func bytes2uint16(b []byte) uint16 {
    p := uint16(0)
    for i := 1;i >= 0;i-- {
        p <<= 8
        p |= uint16(b[i])
    }
    return p
}

func uint322bytes(u uint32) []byte {
    b := make([]byte, 4)
    for i := 0;i < 4;i++ {
         b[i] = (byte)(u)
         u >>= 8
    }
    return b
}

func bytes2uint32(b []byte) uint32 {
    p := uint32(0)
    for i := 3;i >= 0;i-- {
        p <<= 8
        p |= uint32(b[i])
    }
    return p
}

func uint2bytes(p uint) []byte {
    b := make([]byte, 8)
    for i := 0;i < 8;i++ {
         b[i] = (byte)(p)
         p >>= 8
    }
    return b
}

func bytes2uint(b []byte) uint {
    p := uint(0)
    for i := 7;i >= 0;i-- {
        p <<= 8
        p |= uint(b[i])
    }
    return p
}
