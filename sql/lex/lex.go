package lex

import (
	"../../log"
)

type DFA struct {
	accept []bool
	array  [][]int
	class  []int
}

var (
	reservesDfa DFA
	normalDfa   DFA
	ac          map[int]string
	dfaKeyword  *preDfa
	dfaOther    *preDfa
)

type Token struct {
	Kind string
	Raw  []byte
}

func DefineToken(token string, g *nfa) {
	if dfaKeyword == nil {
		fork := g.fork()
		fork.Simplify()
		dfaKeyword = fork.toDfa()
		for i, v := range dfaKeyword.v {
			if v {
				ac[i] = token
			}
		}
	} else {
		fork := g.fork()
		fork.Simplify()
		for _, v := range dfaKeyword.addDfa(fork.toDfa()) {
			ac[v] = token
		}
	}
}

func DefineCommon(token string, g *nfa) {
	if dfaOther == nil {
		fork := g.fork()
		fork.Simplify()
		dfaOther = fork.toDfa()
		for i, v := range dfaOther.v {
			if v {
				ac[i] = token
			}
		}
	} else {
		fork := g.fork()
		fork.Simplify()
		for _, v := range dfaOther.addDfa(fork.toDfa()) {
			ac[v] = token
		}
	}
}

func init() {
	ac = make(map[int]string)
	defineTokens()
	reservesDfa.class, reservesDfa.array, reservesDfa.accept = dfaKeyword.toArray()
	normalDfa.class, normalDfa.array, normalDfa.accept = dfaOther.toArray()
	log.WriteLogSync("sys", "lex module ready")
}

func numberNfa() *nfa {
	nfas := make([]*nfa, 0)
	for i := 0; i < 9; i++ {
		nfas = append(nfas, single(byte(int('0')+i)))
	}
	return or(nfas...)
}

func letterNfa() *nfa {
	nfas := make([]*nfa, 0)
	for i := 0; i < 26; i++ {
		nfas = append(nfas, single(byte(int('a')+i)))
	}
	for i := 0; i < 26; i++ {
		nfas = append(nfas, single(byte(int('A')+i)))
	}
	return or(nfas...)
}

func stringsToken(tokens ...string) *nfa {
	ret := make([]*nfa, 0)
	for _, v := range tokens {
		ret = append(ret, strings([]byte(v)))
	}
	return or(ret...)
}

func Parse(input ByteReader) ([]Token, error) {
	log.WriteLogSync("query", "Parser:"+string(input.data))
	defer log.WriteLogSync("query", "Parser finished")
	t := make([]Token, 0)
	in := ByteReader{
		data: input.data,
	}
	var err error
	for {
		fork := ByteReader{
			data: input.data,
			pos:  in.pos,
		}
		s1, b1, err1 := RunDFA(reservesDfa.class, reservesDfa.array, reservesDfa.accept, fork)
		s2, b2, err2 := RunDFA(normalDfa.class, normalDfa.array, normalDfa.accept, fork)
		var s int
		var b []byte
		if b1 == nil {
			if b2 == nil {
				break
			}
			b = b2
			s = s2
			err = err2
		} else if b2 != nil {
			if len(b2) > len(b1) {
				b = b2
				s = s2
				err = err2
			} else {
				b = b1
				s = s1
				err = err1
			}
		} else {
			b = b1
			s = s1
			err = err1
		}
		token := Token{}
		if k, ok := ac[s]; ok {
			token.Kind = k
		}
		token.Raw = b
		t = append(t, token)
		if err == nil {
			break
		}
		in.pos += len(b)
		if in.pos >= len(in.data) {
			break
		}
	}
	return t, err
}
