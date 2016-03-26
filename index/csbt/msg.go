package csbt

import (
    "fmt"
    "strconv"
)

type Msg struct {
    Info    string
    Time    int64
}

func (m *Msg) Print() {
    fmt.Println(m.Info)
    fmt.Println("With in " + strconv.FormatFloat(float64(m.Time) / 1000000000.0, 'f', -1, 64) + "s.")
}