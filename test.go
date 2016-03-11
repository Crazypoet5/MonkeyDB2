package main 
import (
    "unsafe"
    "fmt"
)

type s struct {
    x int
    y int
}

func main() {
    ss := make([]s, 4)
    ss[2].x = 5
    ss[2].y = 3
    var k interface{}
    k = ss[0]
    fmt.Print(unsafe.Sizeof(k.(s)))
}