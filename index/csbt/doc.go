package csbt

func (tree *DCSBT) Select(key uint) uint {
    return searchCSBTNode(tree.indexHeader.root, key)
}

func (tree *DCSBT) Insert(key, value uint, base uintptr) error {
    return insertToDCSBT(tree, key, value, base)
}