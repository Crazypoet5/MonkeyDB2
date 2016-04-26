package syntax

import (
	"../lex"
)

type TokenReader struct {
	data []lex.Token
	pos  int
}

func (r *TokenReader) Read() lex.Token {
	r.pos++
	if r.pos > len(r.data) {
		return lex.Token{}
	}
	if r.data[r.pos-1].Kind != "split" {
		return r.data[r.pos-1]
	} else {
		return r.Read()
	}
}

func (r *TokenReader) DirectRead() lex.Token {
	r.pos++
	if r.pos > len(r.data) {
		return lex.Token{}
	}
	return r.data[r.pos-1]
}

func (r *TokenReader) Empty() bool {
	return r.pos >= len(r.data)-1
}

func (r *TokenReader) Fork() *TokenReader {
	return &TokenReader{
		data: r.data,
		pos:  r.pos,
	}
}

func (r *TokenReader) Next(n int) {
	for i := 0; i < n; i++ {
		r.Read()
	}
}

func NewTokenReader(data []lex.Token) *TokenReader {
	return &TokenReader{
		data: data,
		pos:  0,
	}
}
