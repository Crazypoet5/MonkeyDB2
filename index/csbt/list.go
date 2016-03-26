package csbt

import (
    "fmt"
)

func (t *DCSBT) ListByLevel() {
    root := t.mb.GetRoot()
    if t.mb.IsLeaf(root) {
        t.PrintLeaf(root)
    } else {
        t.PrintNode(root)
    }
}

func (t *DCSBT) PrintLeaf(leaf uint) {
    fmt.Println("Leaf {")
    keyNum :=t.mb.GetLeafKeyNum(leaf)
    fmt.Println("\tKeyNum:", keyNum)
    fmt.Print("\tKeys: ")
    for i := 0;i < keyNum;i++ {
        fmt.Print(t.mb.GetLeafKey(leaf, i), " ")
    }
    fmt.Print("\n")
    fmt.Println("}")
}

func (t *DCSBT) PrintNode(node uint) {
    fmt.Println("Node {")
    keyNum :=t.mb.GetNodeKeyNum(node)
    fmt.Println("\tKeyNum:", keyNum)
    fmt.Print("\tKeys: ")
    for i := 0;i < keyNum;i++ {
        c := t.mb.GetChild(node, i)
        if t.mb.IsLeaf(c) {
            t.PrintLeaf(c)
        } else {
            t.PrintNode(t.mb.GetChild(node, i))
        }
        fmt.Print(t.mb.GetNodeKey(node, i), " ")
    }
    t.PrintLeaf(t.mb.GetChild(node, keyNum))
    fmt.Print("\n")
    fmt.Println("}")
}