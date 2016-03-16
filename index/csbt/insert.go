package csbt

import (
    "unsafe"
    "errors"
)

func insertToDCSBT(tree *DCSBT, key uint, value uint, base uintptr) (error) {
    l, k, b, p := selectCSBTNode(tree.indexHeader.root, key, nil)
    if b {
        return errors.New("already exist key:" + string(key))
    }
    if l.keyNum < NODE_NUM / 2 {
        for i := l.keyNum;i > k;i-- {
            l.key[i] = l.key[i - 1]
        }
        l.key[k] = key
        l.data[k] = value
        l.base[k] = base
        l.keyNum++
        return nil
    }
    if root, ok := tree.indexHeader.root.(*CSBTLeaf);ok { //Root is leaf, key must in root
        sp1 := tree.allocate(LEAF_SIZE)
        splited1 := (*CSBTLeaf)(unsafe.Pointer(&sp1[0]))
        if unsafe.Sizeof(splited1) != (NODE_NUM / 3) * 24 + 24 {
            panic("serialize RAM address error!")
        }
        sp2 := tree.allocate(LEAF_SIZE)
        splited2 := (*CSBTLeaf)(unsafe.Pointer(&sp2[0]))
        tree.indexHeader.min = splited1
        tree.indexHeader.max = splited2
        splited1.keyNum = NODE_NUM / 6 
        splited1.left = nil
        splited1.right = splited2
        splited2.left = splited1
        splited2.right = nil
        splited2.keyNum = NODE_NUM / 3 - NODE_NUM / 6
        for i := 0;i < NODE_NUM / 6;i++ {
            splited1.key[i] = root.key[i]
            splited1.base[i] = root.base[i]
            splited1.data[i] = root.data[i]
        }
        for i := root.keyNum - 1;i > NODE_NUM/6 - 1;i-- {
            splited1.key[i] = root.key[i]
        }
        r := tree.allocate(NODE_SIZE)
        newRoot := (*CSBTNode)(unsafe.Pointer(&r))
        newRoot.keyNum = 2
        newRoot.key[0] = splited1.key[splited1.keyNum - 1]
        newRoot.key[1] = splited2.key[splited2.keyNum - 1]
        newRoot.child = splited1
        tree.indexHeader.root = newRoot
        splited1.key[splited1.keyNum] = key
        splited1.data[splited1.keyNum] = value
        splited1.base[splited1.keyNum] = base
        return nil
    }
    if p.keyNum != NODE_NUM {
        var last *CSBTLeaf
        last = nil
        var i int
        var first interface{}
        for i = 0;p.key[i] < key;i++ {
            newb := tree.allocate(NODE_SIZE)
            newn := (*CSBTLeaf)(unsafe.Pointer(&newb[0]))
            if i == 0 {
                first = newn
            }
            src := p.children(i).(*CSBTLeaf)
            for j := 0;j < src.keyNum;j++ {
                newn.base[j] = src.base[j]
                newn.data[j] = src.data[j]
                newn.key[j] = src.key[j]
            }
            newn.keyNum = src.keyNum
            newn.left = last
            if last != nil {
                last.right = newn
            }
            
        }
        newb := tree.allocate(NODE_SIZE)
        newn := (*CSBTLeaf)(unsafe.Pointer(&newb[0]))
        newn.keyNum = 1
        newn.key[0] = key
        newn.base[0] = base
        newn.data[0] = value
        newn.left = last
        newn.right = last.right
        last.right = newn
        last = newn
        pos := i
        for ;i < p.keyNum;i++ {
            newb := tree.allocate(NODE_SIZE)
            newn := (*CSBTLeaf)(unsafe.Pointer(&newb[0]))
            src := p.children(i).(*CSBTLeaf)
            for j := 0;j < src.keyNum;j++ {
                newn.base[j] = src.base[j]
                newn.data[j] = src.data[j]
                newn.key[j] = src.key[j]
            }
            newn.keyNum = src.keyNum
            newn.left = last
            if last != nil {
                last.right = newn
            }
        }
        last.right = p.children(p.keyNum - 1).(*CSBTLeaf).right
        for k := p.keyNum;k > pos;k-- {
            p.key[k] = p.key[k - 1]
        }
        p.keyNum++
        p.key[pos] = key
        p.child = first
        return nil
    }
    
    return nil
}
