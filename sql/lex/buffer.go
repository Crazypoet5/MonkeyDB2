package lex

import "errors"

type ByteReader struct {
    data []byte
    pos int
}

func (br *ByteReader) Fork() *ByteReader {
    return &ByteReader {
        data:   br.data,
        pos:    br.pos,
    }
}

func (br *ByteReader) Read() (byte, error) {
    if br.pos >= len(br.data) {
        return 0, errors.New("Empty")
    }
    br.pos++
    return br.data[br.pos - 1], nil
}

func (br *ByteReader) Empty() bool {
    return br.pos >= len(br.data)
}

func NewByteReader(bytes []byte) *ByteReader {
    return &ByteReader {
        data:   bytes,
    }
}