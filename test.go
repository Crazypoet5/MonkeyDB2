package main 
import (
    "fmt"
    "unsafe"
)

type p struct {
    x int
    y int
}

func main() {
    pp := p {
        x:  1,
        y:  2,
    }
    b := (*[]byte)(unsafe.Pointer(&pp))
    fmt.Println(*b)
}