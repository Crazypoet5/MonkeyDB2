package memory

import (
    "testing"
    "unsafe"
)

var I uintptr

func TestCreateImage(t *testing.T) {
    I, err := CreateImage(1024)
    if err != nil {
        t.Error(err)
    }
    (*[5]int)(unsafe.Pointer(I))[2] = 5
    if _, ok := ImageTable[I];!ok {
        t.Error("Not in map")
    }
    ReleaseImage(I)
}

func TestReallocImage(t *testing.T) {
    I, err := CreateImage(1024)
    if err != nil {
        t.Error(err)
    }
    (*[5]int)(unsafe.Pointer(I))[2] = 5
    if _, ok := ImageTable[I];!ok {
        t.Error("Not in map")
    }
    newI, err := ReallocImage(I, 2048)
    if err != nil {
        t.Error(err)
    }
    (*[5]int)(unsafe.Pointer(newI))[2] = 5
    if _, ok := ImageTable[I];ok && I != newI {
        t.Error("Old in map")
    }
    if _, ok := ImageTable[newI];!ok {
        t.Error("New not in map")
    }
    ReleaseImage(newI)
    I = newI
}

func TestReleaseImage(t *testing.T) {
    I, err := CreateImage(1024)
    if err != nil {
        t.Error(err)
    }
    (*[5]int)(unsafe.Pointer(I))[2] = 5
    if _, ok := ImageTable[I];!ok {
        t.Error("Not in map")
    }
    ReleaseImage(I)
    ReleaseImage(I)
    if _, ok := ImageTable[I];ok {
        t.Error("In map")
    }
}