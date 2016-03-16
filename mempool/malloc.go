package mempool

import (
    "strconv"
    "../log"
    "time"
    "reflect"
    "unsafe"
)

// MallocTable record which address was based on which file
var MallocTable = make(map[uintptr]fileImage)


// Please use GetFree instead or you might make an error
func Malloc(size int) []byte{
    defer syncTableToFile()
    datetime := strconv.Itoa(int(time.Now().UnixNano()))
    filename := ".\\image" + datetime
    h := CreateFile(filename, OPEN_ALWAYS)
    hI := CreateFileMapping(h, 0, uint(size), "img" + datetime)
    ip := MapViewOfFile(hI, uint(size))
    //Use the sync File IO appears to make OS refresh map view
    log.WriteLogSync("sys", "MapViewOfFile return:" + strconv.Itoa(int(ip)))
    MallocTable[ip] = fileImage {
        FileHandle: h,
        ImageHandle:hI,
        FileName:   filename,
    }
    var header reflect.SliceHeader
    header.Data = ip
    header.Len = size
    header.Cap = size
    b := *(*[]byte)(unsafe.Pointer(&header))
    for i := 0;i < size;i++ {
        b[i] = 0
    }
    FlushViewOfFile(ip, uint(size))
    //UnmapViewOfFile(ip)
    //CloseHandle(hI)
    CloseHandle(h)
    return b
}

// Please use Release instead or you might make an error
func Free(p []byte) {
    defer syncTableToFile()
    header := (*reflect.SliceHeader)(unsafe.Pointer(&p))
    UnmapViewOfFile(header.Data)
    CloseHandle(MallocTable[header.Data].ImageHandle)
    //CloseHandle(MallocTable[header.Data].fileHandle)
    delete(MallocTable, header.Data)
    header.Len = 0
    header.Cap = 0
    header.Len = 0
}
