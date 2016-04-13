package csbt

import (
	"strconv"

	"../../log"
	"../../memory"
)

func (t *DCSBT) Recovery() {
	for l := t.mb.GetMin(); l != 0; l = t.mb.GetLeafRight(l) {
		keyNum := t.mb.GetLeafKeyNum(l)
		for i := 0; i < keyNum; i++ {
			oldPtr := t.mb.GetLeafValue(l, i)
			newPtr, ok := memory.RecoveryTable[oldPtr]
			if !ok {
				log.WriteLog("err", "Recovery index error: "+strconv.Itoa(int(oldPtr)))
				t.mb.SetLeafValue(l, i, 0)
				continue
			}
			t.mb.SetLeafValue(l, i, newPtr.RawPtr)
		}
	}
}
