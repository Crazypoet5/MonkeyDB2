package syntax

import (
	"testing"

	"../lex"
)

func Test_valueParser(t *testing.T) {
	ts, _ := lex.Parse(*lex.NewByteReader([]byte("123.5")))
	stn, err := valueParser(&TokenReader{data: ts})
	if err != nil {
		t.Error(err)
	}
	stn.Print(1)
}
