package table

type Value []byte

func (v *Value) Len() int {
    return len(*v)
}

func (v *Value) AsInt() int {
    return int(bytes2uint(*v))
}

func (v *Value) AsFloat() float64 {
    return 0.0
}

func (v *Value) AsString() string {
    return string(*v)
}