package csbt

import (
    "../msg"
    "time"
)

func (t *DCSBT) delete(key uint32) *msg.Msg {
    start := time.Now().UnixNano()
    l, _, i, b := t.selecT(t.mb.GetRoot(), 0, key)
    if !b {
        return &msg.Msg {
            Info:   "Not found.",
            Success:    false,
        }
    }
    keyNum := t.mb.GetLeafKeyNum(l)
    for i := i;i < keyNum - 1;i++ {
        t.mb.SetLeafKey(l, i, t.mb.GetLeafKey(l, i + 1))
        t.mb.SetLeafValue(l, i, t.mb.GetLeafValue(l, i + 1))
    }
    t.mb.SetLeafKeyNum(l, keyNum - 1)
    return &msg.Msg {
        Info:   "Delete success.",
        Success:    true,
        Time:   time.Now().UnixNano() - start,
    }
}