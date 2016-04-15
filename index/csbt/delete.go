package csbt

//Delete is only to set the value to 0
func (t *DCSBT) Delete(key uint32) {
	n := t.FindLeaf(key, false)
	l := t.MB.GetChild(n.node, n.offset)
	for i := 0; i < t.MB.GetLeafKeyNum(l); i++ {
		if t.MB.GetLeafKey(l, i) == key {
			t.MB.SetLeafValue(l, i, 0)
		}
	}
}
