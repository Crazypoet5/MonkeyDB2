package syntax

import (
	"strconv"
)

func calculateInt(raw []byte) int {
	i, err := strconv.Atoi(string(raw))
	if err != nil {
		return 0
	}
	return i
}

func calculateFloat(raw []byte) float64 {
	f, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		fI, err := (strconv.Atoi(string(raw)))

		if err != nil {
			return 0.0
		} else {
			f = float64(fI)
			return f
		}
	}
	return f
}
