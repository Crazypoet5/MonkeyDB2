package syntax

import (
	"testing"

	"../lex"
)

func Test_rowParser(t *testing.T) {
	ts, _ := lex.Parse(*lex.NewByteReader([]byte("(123.5, 'a')")))
	stn, err := rowParser(&TokenReader{data: ts})
	if err != nil {
		t.Error(err)
	}
	stn.Print(1)
}
