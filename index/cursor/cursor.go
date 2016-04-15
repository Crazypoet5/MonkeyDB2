package cursor

type Cursor interface {
	Next() Cursor
	Prev() Cursor
	Read() (uintptr, uint)
	Delete()
	Write(uintptr, uint)
}
