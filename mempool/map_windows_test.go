package mempool

import "testing"
import "unsafe"

func TestCreateFile(t *testing.T) {
    h := CreateFile("./a.txt", OPEN_ALWAYS)
    if h == 0 {
        t.Error(0)
    }
    CloseHandle(h)
}

func TestCreateFileMapping(t *testing.T) {
    h := CreateFile("./a.img", OPEN_ALWAYS)
    hI := CreateFileMapping(h, 0, 32768, "img")
    ip := MapViewOfFile(hI, 4)
    p := (*int)(unsafe.Pointer(ip))
    *p = 5
    FlushViewOfFile(ip, 4)
    UnmapViewOfFile(ip)
    CloseHandle(hI)
    CloseHandle(h)
    h = CreateFile("./a.img", OPEN_ALWAYS)
    hI = CreateFileMapping(h, 0, 32768, "img")
    ip = MapViewOfFile(hI, 4)
    p = (*int)(unsafe.Pointer(ip))
    if *p != 5 {
        t.Error("Failed to check img file and RAM.\n")
        CloseHandle(hI)
        CloseHandle(h)
    }
    CloseHandle(hI)
    CloseHandle(h)
}