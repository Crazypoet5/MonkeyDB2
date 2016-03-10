package syntax

import (
    "testing"
    "../lex"
)

func Test_filedParser(t *testing.T) {
    ts, _ := lex.Parse(*lex.NewByteReader([]byte("abc.abc.efd")))
    stn, err := filedParser(&TokenReader{data: ts})
    if err != nil {
        t.Error(err)
    }
    stn.Print(1)
}