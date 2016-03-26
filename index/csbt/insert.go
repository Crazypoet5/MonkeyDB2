package csbt

import (
    "time"
)

func (t *DCSBT) Insert(key uint32, value uintptr) *Msg {
    start := time.Now().UnixNano()
    l, p, i, b := t.selecT(t.mb.GetRoot(), 0, key)
    if b {
        return &Msg {
            Info:   "Data already exist.",
            Time:   time.Now().UnixNano() - start,
        }
    }
    if t.mb.GetLeafKeyNum(l) < 3 {
        t.insertToFreeLeaf(l, i, key, value)
    } else {
        t.splitFullLeafAndInsert(l, i, p, key, value)
    }
    return &Msg {
        Info:   "Insert ok.",
        Time:   time.Now().UnixNano() - start,
    }
}

func (t *DCSBT) insertToFreeLeaf(leaf uint, index int, key uint32, value uintptr) {
    keyNum := t.mb.GetLeafKeyNum(leaf)
    for i := keyNum;i > index + 1;i-- {
        k := t.mb.GetLeafKey(leaf, i - 1)
        v := t.mb.GetLeafValue(leaf, i - 1)
        t.mb.SetLeafKey(leaf, i, k)
        t.mb.SetLeafValue(leaf, i, v)
    }
    t.mb.SetLeafKey(leaf, index + 1, key)
    t.mb.SetLeafValue(leaf, index + 1, value)
    t.mb.SetLeafKeyNum(leaf, keyNum + 1)
}

func (t *DCSBT) splitFullLeafAndInsert(leaf uint, index int, p uint, key uint32, value uintptr) {
    if p == 0 {     //Leaf node is root
        p = t.mb.NewNodes(2)
        t.mb.SetRoot(p)
        t.mb.SetChildren(p, leaf)
        t.mb.SetNodeKeyNum(p, 1)
    }
    keyNum := t.mb.GetNodeKeyNum(p)
    if keyNum < 13 {
        f := t.mb.GetChild(p, 0)
        left := t.mb.GetLeafLeft(f)
        right := t.mb.GetLeafRight(f)
        g0 := t.mb.NewLeaves(keyNum + 2)
        prevLeaf, _ := t.mb.Read(f, leaf - f)
        t.mb.Write(g0, prevLeaf)
        new1 := g0 + uint(len(prevLeaf))
        new2 := new1 + 64
        afterLeaf, _ := t.mb.Read(leaf + 64, f + 64 * uint(keyNum + 1) - leaf - 64)
        t.mb.Write(g0 + uint(len(prevLeaf)) + 128, afterLeaf)
        k0 := t.mb.GetLeafKey(leaf, 0)
        k1 := t.mb.GetLeafKey(leaf, 1)
        k2 := t.mb.GetLeafKey(leaf, 2)
        v0 := t.mb.GetLeafValue(leaf, 0)
        v1 := t.mb.GetLeafValue(leaf, 2)
        v2 := t.mb.GetLeafValue(leaf, 2)
        switch(index) {
            case -1:
            t.mb.SetLeafKey(new1, 0, key)
            t.mb.SetLeafValue(new1, 0, value)
            t.mb.SetLeafKey(new1, 1, k0)
            t.mb.SetLeafValue(new1, 1, v0)
            t.mb.SetLeafKeyNum(new1, 2)
            t.mb.SetLeafKey(new1, 0, k1)
            t.mb.SetLeafValue(new2, 0, v1)
            t.mb.SetLeafKey(new1, 0, k2)
            t.mb.SetLeafValue(new2, 1, v2)
            t.mb.SetLeafKeyNum(new2, 2)
            case 0:
            t.mb.SetLeafKey(new1, 0, k0)
            t.mb.SetLeafValue(new1, 0, v0)
            t.mb.SetLeafKey(new1, 1, key)
            t.mb.SetLeafValue(new1, 1, value)
            t.mb.SetLeafKeyNum(new1, 2)
            t.mb.SetLeafKey(new1, 0, k1)
            t.mb.SetLeafValue(new2, 0, v1)
            t.mb.SetLeafKey(new1, 0, k2)
            t.mb.SetLeafValue(new2, 1, v2)
            t.mb.SetLeafKeyNum(new2, 2)
            case 1:
            t.mb.SetLeafKey(new1, 0, k0)
            t.mb.SetLeafValue(new1, 0, v0)
            t.mb.SetLeafKey(new1, 1, k1)
            t.mb.SetLeafValue(new1, 1, v1)
            t.mb.SetLeafKeyNum(new1, 2)
            t.mb.SetLeafKey(new1, 0, key)
            t.mb.SetLeafValue(new2, 0, value)
            t.mb.SetLeafKey(new1, 0, k2)
            t.mb.SetLeafValue(new2, 1, v2)
            t.mb.SetLeafKeyNum(new2, 2)
            case 2:
            t.mb.SetLeafKey(new1, 0, k0)
            t.mb.SetLeafValue(new1, 0, v0)
            t.mb.SetLeafKey(new1, 1, k1)
            t.mb.SetLeafValue(new1, 1, v1)
            t.mb.SetLeafKeyNum(new1, 2)
            t.mb.SetLeafKey(new1, 0, key)
            t.mb.SetLeafValue(new2, 0, value)
            t.mb.SetLeafKey(new1, 0, k2)
            t.mb.SetLeafValue(new2, 1, v2)
            t.mb.SetLeafKeyNum(new2, 2)
        }
        //Repair left, right and isLeaf
        for p := g0;p < g0 + 64 * uint(keyNum + 2);p += 64 {
            t.mb.SetLeaf(p)
            t.mb.SetLeafLeft(p, p - 64)
            t.mb.SetLeafRight(p, p + 64)
        }
        t.mb.SetLeafLeft(g0, left)
        t.mb.SetLeafRight(g0 + 64 * uint(keyNum), right)
        t.mb.SetNodeKeyNum(p, keyNum + 1)
        t.mb.SetChildren(p, g0)
        if p == t.mb.GetRoot() {
            t.mb.SetMin(g0)
            t.mb.SetMax(g0 + 64 * uint(keyNum))
        }
        return
    }
    //TODO
    pp := t.findParent(p)
    resplit:
    if pp == 0 {    //p is root
        pp = t.mb.NewNodes(1)
        t.mb.SetRoot(pp)
        t.mb.SetChildren(pp, p)
        t.mb.SetNodeKeyNum(pp, 1)
    }
    keyNum = t.mb.GetNodeKeyNum(pp)
    if keyNum < 13 {
        t.splitNode(p, pp)
        t.Insert(key, value)
    } else {
        p = pp
        pp = t.findParent(p)
        goto resplit
    }
}

//p is the parent of node and not full
func (t *DCSBT) splitNode(node, p uint) {
    f := t.mb.GetChild(p, 0)
    keyNum := t.mb.GetNodeKeyNum(p)
    g0 := t.mb.NewNodes(keyNum + 2)
    prevNode, _ := t.mb.Read(f, node - f)
    t.mb.Write(g0, prevNode)
    new1 := g0 + uint(len(prevNode))
    new2 := new1 + 64
    afterNode, _ := t.mb.Read(node + 64, f + 64 * uint(keyNum + 1) - node - 64)
    t.mb.Write(g0 + uint(len(prevNode)) + 128, afterNode)
    t.mb.SetNodeKeyNum(new1, 7)
    t.mb.SetChildren(new1, t.mb.GetChild(node, 0))
    t.mb.SetNodeKeyNum(new2, 6)
    for i := 0;i < 6;i++ {
        t.mb.SetNodeKey(new2, i, t.mb.GetNodeKey(node, 7 + i))
    }
    t.mb.SetChildren(new2, t.mb.GetChild(node, 7))
    indexOfNode := int((node - f) / 64)
    for i := keyNum;i > indexOfNode + 1;i-- {
        t.mb.SetNodeKey(p, i, t.mb.GetNodeKey(p, i-1))
    }
    t.mb.SetNodeKey(p, indexOfNode + 1, t.mb.GetNodeKey(node, 7))
    t.mb.SetChildren(p, g0)
    t.mb.SetNodeKeyNum(p, keyNum + 1)
}

func (t *DCSBT) findParent(node uint) uint {
    key := t.mb.GetNodeKey(node, 0)
    p := t.mb.GetRoot()
    for {
        if t.mb.IsLeaf(p) {
            return 0
        }
        child := t.mb.GetChild(p, 0)
        keyNum := t.mb.GetNodeKeyNum(p)
        if child <= node && child + 64 * uint(keyNum) > node {
            return p
        }
        i := 0
        for ;i < keyNum;i++ {
            k := t.mb.GetNodeKey(p, i)
            if k > key {
                break
            }
        }
        p = t.mb.GetChild(p, i + 1)
    }
}