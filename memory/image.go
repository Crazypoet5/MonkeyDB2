package memory

//#include <stdlib.h>
import "C"

import (
    "strconv"
    "time"
    "../common"
    "unsafe"
    "os"
)

var commonBuffer = make([]byte, 1024 * 1024)    // To clear the file quickly

// CreateImage creates a image file and returns the address
func CreateImage(size int) (ip uintptr, err error) {
    filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(int(time.Now().UnixNano()))
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

// ReallocImage creates a new bigger image file and returns the new address, but not copy
func ReallocImage(ip uintptr, size int) (uintptr, error) {
    filename := common.COMMON_DIR + "\\image\\" + strconv.Itoa(int(time.Now().UnixNano()))
    C.free(unsafe.Pointer(ip))
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
    ImageTable[ipNew] = filename
    delete(ImageTable, ip)
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