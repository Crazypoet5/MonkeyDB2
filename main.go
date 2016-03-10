package main

import (
    "syscall"
    "unsafe"
)

var (
    user32, _ = syscall.LoadLibrary("user32.dll")
)

var (
    messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
)

func StrPtr(str string) uintptr {
    return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str)))
}

func MessageBox(caption, text string, style uintptr) {
    _, _, callErr := syscall.Syscall9(messageBox, 4, 0, StrPtr(text), StrPtr(caption), style, 0, 0, 0, 0, 0)
    if callErr != 0 {
        panic(callErr)
    }
}

func main() {
    MessageBox("Caption", "Alert", 0)
}