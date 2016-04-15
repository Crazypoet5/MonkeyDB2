package csbt

import (
	"fmt"
)

func (t *DCSBT) ListByLevel() {
	root := t.MB.GetRoot()
	if t.MB.IsLeaf(root) {
		t.PrintLeaf(root, "")
	} else {
		t.PrintNode(root, "")
	}
}

func (t *DCSBT) PrintLeaf(leaf uint, tabs string) {
	fmt.Println(tabs + "Leaf {")
	keyNum := t.MB.GetLeafKeyNum(leaf)
	fmt.Println(tabs+"\tKeyNum:", keyNum)
	fmt.Print(tabs + "\tKeys: ")
	for i := 0; i < keyNum; i++ {
		fmt.Print(t.MB.GetLeafKey(leaf, i), " ")
	}
	fmt.Print("\n")
	fmt.Println(tabs + "}")
}

func (t *DCSBT) PrintNode(node uint, tabs string) {
	fmt.Println(tabs + "Node {")
	keyNum := t.MB.GetNodeKeyNum(node)
	fmt.Println(tabs+"\tKeyNum:", keyNum)
	fmt.Println(tabs + "\tKeys: ")
	for i := 0; i < keyNum; i++ {
		c := t.MB.GetChild(node, i)
		if t.MB.IsLeaf(c) {
			t.PrintLeaf(c, tabs+"\t")
		} else {
			t.PrintNode(t.MB.GetChild(node, i), tabs+"\t")
		}
		fmt.Println(tabs+"\t", t.MB.GetNodeKey(node, i))
	}
	t.PrintLeaf(t.MB.GetChild(node, keyNum), tabs+"\t")
	fmt.Print("\n")
	fmt.Println(tabs + "}")
}
