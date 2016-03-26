package csbt

import (
)

// Return index is less than key
func (t *DCSBT) selecT(start, parent uint, key uint32) (uint, uint, int, bool) {
    if t.mb.IsLeaf(start) {
        keyNum := t.mb.GetLeafKeyNum(start)
        if keyNum == 0 {
            return start, parent, -1, false
        }
        for i := 0;i < keyNum;i++ {
            k := t.mb.GetLeafKey(start, i)
            if k == key {
                return start, parent, i, true
            }
            if k > key {
                return start, parent, i - 1, false
            }
        }
        return start, parent, keyNum - 1, false
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

func (t *DCSBT) SelectByConfition(condition func(uint32) bool) []uintptr {
    ret := make([]uintptr, 0)
    for l := t.mb.GetMin();l != 0;l = t.mb.GetLeafRight(l) {
        keyNum := t.mb.GetLeafKeyNum(l)
        for i := 0;i < keyNum;i++ {
            key := t.mb.GetLeafKey(l, i)
            if condition(key) {
                value := t.mb.GetLeafValue(l, i)
                ret = append(ret, value)
            }
        }
    }
    return ret
}

func abs(n int) int {
    if n > 0 {
        return n
    }
    return -n
}

func (t *DCSBT) Select(key uint32, n int) []uintptr {
    ret := make([]uintptr, 0)
    l, _, i, b := t.selecT(t.mb.GetRoot(), 0, key)
    if !b {
        return ret
    }
    sum := 0
    next := func(l uint) uint {
        if n > 0 {
            return t.mb.GetLeafRight(l)
        }
        return t.mb.GetLeafLeft(l)
    }
    for e := l;e != 0 && sum < abs(n);e = next(e) {
        j := 0
        keyNum := t.mb.GetLeafKeyNum(e)
        if e == l {
            j = i
        }
        for ;j < keyNum;j++ {
            ret = append(ret, t.mb.GetLeafValue(e, j))
            sum++
        }
    }
    return ret[:abs(n)]
}