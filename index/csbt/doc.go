package csbt

func (tree *DCSBT) Select(key uint, n int) (r []uint) {
    r = make([]uint, 0)
    l, i, b, _ := selectCSBTNode(tree.indexHeader.root, key, nil)
    if !b {
        return
    }
    sum := 0
    for iterator := l;iterator != nil;iterator = iterator.right {
        for j := i;j < iterator.keyNum;j++ {
            r = append(r, iterator.data[j])
            sum++
            if n != 0 && sum >= n {
                return
            }
        }
        i = 0
    }
    return
}

func (tree *DCSBT) Insert(key, value uint, base uintptr) error {
    return insertToDCSBT(tree, key, value, base)
}

func (tree *DCSBT) InsertBat(key, value []uint, base []uintptr) error {
    return nil
}