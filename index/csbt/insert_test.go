package csbt

import (
    "testing"
    "fmt"
)

func TestInsert(t *testing.T) {
    fmt.Println("S")
    tree := NewDCSBT()
    tree.Insert(3, 123)
    tree.Insert(4, 123)
    tree.Insert(1, 123)
    tree.ListByLevel()
    tree.Insert(5, 123)
    tree.ListByLevel()
    tree.Insert(6, 123)
    tree.ListByLevel()
    tree.Insert(8, 123)
    tree.ListByLevel()
    tree.Insert(4, 123)
    tree.ListByLevel()
}