package common

import (
	"fmt"
)

var PRINT2 = false

func Print2(str string) {
	if !PRINT2 {
		return
	}
	fmt.Println("[PERFORMANCE TEST]:", str)
}
