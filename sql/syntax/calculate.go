package syntax

import (
    "strconv"
)

func calculateInt(raw []byte) int {
    i, _ := strconv.Atoi(string(raw))
    return i
}

func calculateFloat(raw []byte) float64 {
    f, _ := strconv.ParseFloat(string(raw), 64)
    return f
}