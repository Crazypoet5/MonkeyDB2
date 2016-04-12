package lex

var NFA *statedNfa
var DFA *statedDfa

type Token struct {
	Kind string
	Raw  []byte
}

func init() {
	NFA = NewBasic()
	//	defineTokens()
	DFA = NFA.ToDFA()
}

func Parse(input ByteReader) ([]Token, error) {
	return DFA.Parse(input)
}
