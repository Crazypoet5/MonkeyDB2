package csbt

import (
)

func (t *DCSBT) selecT(start, parent uint, key uint32) (uint, uint, int, bool) {
    if t.mb.IsLeaf(start) {
        keyNum := t.mb.GetLeafKeyNum(start)
        for i := 0;i < keyNum;i++ {
            k := t.mb.GetLeafKey(start, i)
            if k == key {
                return start, parent, i, true
            }
            if k > key {
                return start, parent, i - 1, false
            }
        }
        return start, parent, -1, false
    }
    keyNum := t.mb.GetNodeKeyNum(start)
    for i := 0;i < keyNum;i++ {
        k := t.mb.GetNodeKey(start, i)
        if k == key {
            child := t.mb.GetChild(start, i + 1)
            return t.selecT(child, start, key)
        }
        if k > key {
            child := t.mb.GetChild(start, i)
            return t.selecT(child, start, key)
        }
    }
    child := t.mb.GetChild(start, keyNum)
    return t.selecT(child, start, key)
}