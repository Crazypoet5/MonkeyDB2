package memory

import (
	"unsafe"
	"reflect"
    "fmt"
)

var CopyTable = make(map[*DataBlock]*DataBlock) //Src -> Dst

func (b *DataBlock) read(offset, size int) ([]byte, error) {
    if offset + size > b.Size {
        return nil, OUT_OF_SIZE
    }
    var header reflect.SliceHeader
    header.Data = uintptr(b.RawPtr + uintptr(offset))
    header.Len = size
    header.Cap = size
    return *(*[]byte)(unsafe.Pointer(&header)), nil
}

func (b *DataBlock) Read(offset, size int) ([]byte, error) {
    b.RWMutex.RLock()
    defer b.RWMutex.RUnlock()
    return b.read(offset, size)
}

func (b *DataBlock) write(offset int, data []byte) (int, error) {
    var header reflect.SliceHeader
    size := len(data)
    header.Data = uintptr(b.RawPtr + uintptr(offset))
    header.Len = size
    header.Cap = size
    d := *(*[]byte)(unsafe.Pointer(&header))
    var n int
    if offset + size > b.Size {
        n = b.Size - offset
    } else {
        n = size
    }
    copy(d, data[:n])
    return n, nil
}

func (b *DataBlock) Write(offset int, data []byte) (int, error) {
    b.RWMutex.Lock()
    defer b.RWMutex.Unlock()
    var copies *DataBlock
    copies, ok := CopyTable[b]
    if !ok {
        return b.write(offset, data)
    }
    copies.Write(offset, data)
    return b.Write(offset, data)
}

func Copy(dst, src *DataBlock) (int, error) {
    CopyTable[src] = dst
    data, err := src.Read(0, src.Size)
    if err != nil {
        return 0, err
    }
    delete(CopyTable, src)
    fmt.Println(data)
    return dst.Write(0, data)
}