package csbt

import "testing"

func expect(exp int64, real int64, t *testing.T) {
    if real != exp {
        t.Error("Expected:", exp, "Get:", real)
    }
}

func Test_searchCSBTNode(t *testing.T) {
    tree := NewDCSBT()
    searchCSBTNode(tree.indexHeader.root, 0)
    insertToDCSBT(tree, 156, 188, 165)
    expect(int64(188), int64(searchCSBTNode(tree.indexHeader.root, 156)), t)
}