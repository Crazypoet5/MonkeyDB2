package csbt

import "testing"
import "fmt"

func expect(exp int, real int, t *testing.T) {
    if real != exp {
        t.Error("Expected:", exp, "Get:", real)
    }
}

func Test_searchCSBTNode(t *testing.T) {
    tree := NewDCSBT()
    searchCSBTNode(tree.indexHeader.root, 0)
    insertToDCSBT(tree, 156, 188, 165)
    fmt.Println(searchCSBTNode(tree.indexHeader.root, 156))
}