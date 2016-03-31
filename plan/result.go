package plan

import (
	"time"
)

type Result struct {
	startTime    int64
	UsedTime     int64
	AffectedRows int
}

func NewResult() *Result {
	return &Result{
		startTime: time.Now().UnixNano(),
	}
}

func (r *Result) SetResult(affectedRows int) {
	r.AffectedRows = affectedRows
	r.UsedTime = time.Now().UnixNano() - r.startTime
}
