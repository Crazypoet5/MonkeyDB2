package syntax

import "../lex"

type TokenReader struct {
    data []lex.Token
    pos int
}

func (r *TokenReader) Read() lex.Token {
    r.pos++
    if r.pos > len(r.data) {
        return lex.Token{}
    }
    if r.data[r.pos - 1].Kind != "split" {
        return r.data[r.pos - 1]
    } else {
        return r.Read()
    }
}

func (r *TokenReader) Empty() bool {
    return r.pos >= len(r.data)
}

func (r *TokenReader) Fork() *TokenReader {
    return &TokenReader {
        data:   r.data,
        pos:    r.pos,
    }
}