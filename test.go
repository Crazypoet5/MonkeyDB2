package main 
import (
    "fmt"
    "bytes"
)

func main() {
    b := []byte{'a', 'b'}
    buff := make([]byte, 1)
    input := bytes.NewReader(b)
    input2 := input
    input.Read(buff)
    fmt.Println(buff)
    input2.Read(buff)
    fmt.Println(buff)
}