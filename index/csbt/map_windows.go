package csbt

import (
    "syscall"
    "unsafe"
)

var (
    kernel32, _ = syscall.LoadLibrary("kernel32.dll")
)

var (
    createFileMapping, _ = syscall.GetProcAddress(kernel32, "CreateFileMappingW")
    mapViewOfFile, _ = syscall.GetProcAddress(kernel32, "MapViewOfFile")
    createFile, _ = syscall.GetProcAddress(kernel32, "CreateFileW")
    closeHandle, _ = syscall.GetProcAddress(kernel32, "CloseHandle")
    flushViewOfFile, _ = syscall.GetProcAddress(kernel32, "FlushViewOfFile")
    unmapViewOfFile, _ = syscall.GetProcAddress(kernel32, "UnmapViewOfFile")
)

const (
    //CreateFile
    GENERIC_READ = 0x80000000
    GENERIC_WRITE = 0x40000000
    CREATE_ALWAYS = 2
    CREATE_NEW = 1  //Exist will fail
    OPEN_ALWAYS = 4
    OPEN_EXISTING = 3
    TRUNCATE_EXISTING = 5
    FILE_ATTRIBUTE_ARCHIVE = 32
    FILE_ATTRIBUTE_ENCRYPTED = 16384
    FILE_ATTRIBUTE_HIDDEN = 2
    FILE_ATTRIBUTE_NORMAL = 128
    FILE_ATTRIBUTE_OFFLINE = 4096
    FILE_ATTRIBUTE_READONLY = 1
    FILE_ATTRIBUTE_SYSTEM = 4
    FILE_ATTRIBUTE_TEMPORARY = 256
    //CreateFileMapping
    PAGE_READWRITE = 4
    //MapViewOfFile
    FILE_MAP_WRITE = 2
)

func CreateFileMapping(hFile syscall.Handle, dwMaxiumSizeHigh, dwMaxiumSizeLow uint, name string) syscall.Handle {
    r, _, _ := syscall.Syscall6(createFileMapping, 6, uintptr(hFile), 0, uintptr(PAGE_READWRITE), uintptr(dwMaxiumSizeHigh), uintptr(dwMaxiumSizeLow), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))))
    return syscall.Handle(r)
}

func CreateFile(fileName string, dwCreationDisposition uint) syscall.Handle {
    r, _, _ := syscall.Syscall9(createFile, 7, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(fileName))), uintptr(GENERIC_READ | GENERIC_WRITE), 0, 0, uintptr(dwCreationDisposition), uintptr(FILE_ATTRIBUTE_NORMAL), 0, 0, 0)
    return syscall.Handle(r)
}

func CloseHandle(handle syscall.Handle) {
    syscall.Syscall(closeHandle, 1, uintptr(handle), 0, 0)
}

func MapViewOfFile(hFileMappingObject syscall.Handle, dwNumberOfBytesToMap uint) uintptr {
    r, _, _ := syscall.Syscall6(mapViewOfFile, 5, uintptr(hFileMappingObject), uintptr(FILE_MAP_WRITE), 0, 0, uintptr(dwNumberOfBytesToMap), 0)
    return r
}

func FlushViewOfFile(lpBaseAddress uintptr, dwNumberOfBytesToFlush uint) {
    syscall.Syscall(flushViewOfFile, 2, lpBaseAddress, uintptr(dwNumberOfBytesToFlush), 0)
}

func UnmapViewOfFile(lpBaseAddress uintptr) {
    syscall.Syscall(unmapViewOfFile, 1, lpBaseAddress, 0, 0)
}