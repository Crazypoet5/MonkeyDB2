package recovery

import (
	"time"
)

const (
	FRAME_TIME = 100 * 1000000
)

var LastFrame = time.Now().UnixNano()

func RestoreFrame() {
	defer func() {
		LastFrame = time.Now().UnixNano()
	}()
	if Restoring == 1 {
		return
	}
	if time.Now().UnixNano()-LastFrame > FRAME_TIME {
		Restore()
	}
}
