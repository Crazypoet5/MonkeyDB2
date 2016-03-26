package csbt

import (
    "fmt"
)

func (t *DCSBT) ListByLevel() {
    root := t.mb.GetRoot()
    if t.mb.IsLeaf(root) {
        t.PrintLeaf(root, "")
    } else {
        t.PrintNode(root, "")
    }
}

func (t *DCSBT) PrintLeaf(leaf uint, tabs string) {
    fmt.Println(tabs + "Leaf {")
    keyNum :=t.mb.GetLeafKeyNum(leaf)
    fmt.Println(tabs + "\tKeyNum:", keyNum)
    fmt.Print(tabs + "\tKeys: ")
    for i := 0;i < keyNum;i++ {
        fmt.Print(t.mb.GetLeafKey(leaf, i), " ")
    }
    fmt.Print("\n")
    fmt.Println(tabs + "}")
}

func (t *DCSBT) PrintNode(node uint, tabs string) {
    fmt.Println(tabs + "Node {")
    keyNum :=t.mb.GetNodeKeyNum(node)
    fmt.Println(tabs + "\tKeyNum:", keyNum)
    fmt.Println(tabs + "\tKeys: ")
    for i := 0;i < keyNum;i++ {
        c := t.mb.GetChild(node, i)
        if t.mb.IsLeaf(c) {
            t.PrintLeaf(c, tabs + "\t")
        } else {
            t.PrintNode(t.mb.GetChild(node, i), tabs + "\t")
        }
        fmt.Println(tabs + "\t", t.mb.GetNodeKey(node, i))
    }
    t.PrintLeaf(t.mb.GetChild(node, keyNum), tabs + "\t")
    fmt.Print("\n")
    fmt.Println(tabs + "}")
}