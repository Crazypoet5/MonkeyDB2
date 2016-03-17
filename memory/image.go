package memory

//#include <stdlib.h>
import "C"

import (
    "strconv"
    "../common"
    "unsafe"
    "os"
)

var commonBuffer = make([]byte, 1024 * 1024)    // To clear the file quickly
var count = 0   // Windows only support 100ns level

// CreateImage creates a image file and returns the address
func CreateImage(size int) (ip uintptr, err error) {
    filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(count)
    count++
    ip = uintptr(C.malloc(C.size_t(size)))
    file, err := os.Create(filename)
    defer file.Close()
    for i := size;i > 0;i -= 1024 * 1024 {
        if i < 1024 * 1024 {
            file.Write(commonBuffer[:i])
        } else {
            file.Write(commonBuffer)
        }
    }
    ImageTable[ip] = filename
    return
}

// ReallocImage creates a new bigger image file and returns the new address,
// TODO: WriteBarrier to thread-safe copy
func ReallocImage(ip uintptr, size int) (uintptr, error) {
    filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(count)
    count++
    os.Remove(ImageTable[ip])
    C.free(unsafe.Pointer(ip))  // Malloc may return a address the same as prev
    ipNew := uintptr(C.malloc(C.size_t(size)))
    file, err := os.Create(filename)
    defer file.Close()
    if err != nil {
        return 0, err
    }
    for i := size;i > 0;i -= 1024 * 1024 {
        if i < 1024 * 1024 {
            file.Write(commonBuffer[:i])
        } else {
            file.Write(commonBuffer)
        }
    }
    delete(ImageTable, ip)
    ImageTable[ipNew] = filename
    return ipNew, nil
}

// ReleaseImage release the image and delete its related file
func ReleaseImage(ip uintptr) (err error) {
    _, ok := ImageTable[ip]
    if !ok {
        err = NOT_FOUND_ADDRESS
        return
    }
    err = os.Remove(ImageTable[ip])
    C.free(unsafe.Pointer(ip))
    delete(ImageTable, ip)
    return
}