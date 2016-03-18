package memory

import (
    "testing"
    "unsafe"
)

var I *DataBlock

func TestCreateImage(t *testing.T) {
    I, err := CreateImage(1024)
    if err != nil {
        t.Error(err)
    }
    I.Write(0, []byte{5})
    if _, ok := ImageTable[I.RawPtr];!ok {
        t.Error("Not in map")
    }
    if DataBlockList.Len() != 1 {
        t.Error("Not in list")
    }
    ReleaseImage(I)
}

func TestReallocImage(t *testing.T) {
    I, err := CreateImage(1024)
    if err != nil {
        t.Error(err)
    }
    I.Write(0, []byte{5})
    if _, ok := ImageTable[I.RawPtr];!ok {
        t.Error("Not in map")
    }
    newI, err := ReallocImage(I, 2048)
    if err != nil {
        t.Error(err)
    }
    r, err := newI.Read(0, 1)
    if err != nil {
        t.Error(err)
    }
    if r[0] != 5 {
        t.Error("Not copy in realloc")
    }
    if _, ok := ImageTable[I.RawPtr];ok && I.RawPtr != newI.RawPtr {
        t.Error("Old in map")
    }
    if _, ok := ImageTable[newI.RawPtr];!ok {
        t.Error("New not in map")
    }
    if DataBlockList.Len() != 1 {
        t.Error("Not in list")
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
    if _, ok := ImageTable[I.RawPtr];!ok {
        t.Error("Not in map")
    }
    ReleaseImage(I)
    if DataBlockList.Len() != 0 {
        t.Error("Released in list")
    }
    if _, ok := ImageTable[I.RawPtr];ok {
        t.Error("In map")
    }
}