package csbt

import (
	"strconv"

	"../../log"
	"../../memory"
)

func (t *DCSBT) Recovery() {
	for l := t.MB.GetMin(); l != 0; l = t.MB.GetLeafRight(l) {
		keyNum := t.MB.GetLeafKeyNum(l)
		for i := 0; i < keyNum; i++ {
			oldPtr := t.MB.GetLeafValue(l, i)
			p := uint(oldPtr) >> 24
			offset := uint(oldPtr) & 0x0000000000ffffff
			newPtr, ok := memory.RecoveryTable[uintptr(p)]
			if !ok {
				log.WriteLog("err", "Recovery index error: "+strconv.Itoa(int(oldPtr)))
				t.MB.SetLeafValue(l, i, 0)
				continue
			}
			newV := uint(newPtr.RawPtr) << 24
			newV |= offset
			t.MB.SetLeafValue(l, i, uintptr(newV))
		}
	}
}
