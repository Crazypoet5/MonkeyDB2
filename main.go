package main

import (
    "fmt"
    "unsafe"
)

type a struct {
    i int
    k int
}

type b struct {
    a
}

func main() {
    x := &a {
        i:  5,
    }
    y := &b {
        a:  *x,
    }
    y.a.i = 5
    fmt.Println(x.i)
    fmt.Println(unsafe.Sizeof(y))
}