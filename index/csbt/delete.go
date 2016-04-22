package csbt

//Delete is only to set the value to 0
func (t *DCSBT) Delete(key uint32) {
	c := t.Select(key)
	if c != nil {
		c.Write(0, 0)
	}
}
