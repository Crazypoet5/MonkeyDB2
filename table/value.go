package table

type Value []byte

func (v *Value) Len() {
    return len(Value)
}

func (v *Value) AsInt() int {
    return bytes2int(v)
}

func (v *Value) AsFloat() float64 {
    return 0.0
}

func (v *Value) AsString() string {
    return string(v)
}