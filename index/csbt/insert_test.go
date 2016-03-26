package csbt

import (
    "testing"
    "fmt"
)

func TestInsert(t *testing.T) {
    fmt.Println("S")
    tree := NewDCSBT()
    tree.Insert(3, 123)
    tree.ListByLevel()
    tree.Insert(4, 123)
    tree.ListByLevel()
    tree.Insert(1, 123)
    tree.ListByLevel()
}