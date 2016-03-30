package lex

import (
    "../../log"
)

var (
    accept []bool
    array [][]int
    class []int
    ac map[int]string
    dfa *preDfa
)

type Token struct {
    Kind string
    Raw []byte
}

func DefineToken(token string, g *nfa) {
    if len(ac) == 0 {
        fork := g.fork()
        fork.Simplify()
        dfa = fork.toDfa()
        for i, v := range dfa.v {
            if v {
                ac[i] = token
            }
        }
    } else {
        fork := g.fork()
        fork.Simplify()
        for _, v := range dfa.addDfa(fork.toDfa()) {
            ac[v] = token
        }
    }
}

func init() {
    ac = make(map[int]string)
    defineTokens()
    class, array, accept = dfa.toArray()
    log.WriteLogSync("sys", "lex module ready")
}

func numberNfa() *nfa {
    nfas := make([]*nfa, 0)
    for i := 0;i < 9;i++ {
        nfas = append(nfas, single(byte(int('0') + i)))
    }
    return or(nfas...)
}

func letterNfa() *nfa {
    nfas := make([]*nfa, 0)
    for i := 0;i < 26;i++ {
        nfas = append(nfas, single(byte(int('a') + i)))
    }
    for i := 0;i < 26;i++ {
        nfas = append(nfas, single(byte(int('A') + i)))
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

func defineTokens() {
    log.WriteLogSync("sys", "Start making DFA")
    keyword := stringsToken("select", "from", "where", "update", "delete", "create", "insert", "into", "table", "order", "by")
    logical := stringsToken("and", "or", "not")
    structs := stringsToken("(", ")", ";", ",", ".")
    split := stringsToken(" ", "\t", "\n")
    relations := stringsToken(">", "<", ">=", "<=", "=", "<>")
    types := stringsToken("int", "float", "varchar", "object", "array")
    DefineToken("floatval", links(repeat(numberNfa()), single('.'), chosable(repeat(numberNfa()))))
    DefineToken("intval", repeat(numberNfa()))
    DefineToken("identical", links(repeat(letterNfa()), chosable(repeat(or(numberNfa(), letterNfa())))))
    DefineToken("keyword", keyword)
    DefineToken("logical", logical)
    DefineToken("structs", structs)
    DefineToken("split", split)
    DefineToken("relations", relations)
    DefineToken("types", types)
    DefineToken("unReference", single('`'))
    DefineToken("reference", single('\''))
    log.WriteLogSync("sys", "DFA prepared")
}

func Parse(input ByteReader) ([]Token, error) {
    log.WriteLogSync("query", "Parser:" + string(input.data))
    defer log.WriteLogSync("query", "Parser finished")
    t := make([]Token, 0)
    in := ByteReader {
        data:   input.data,
    }
   var err error
    for {
        fork := ByteReader {
            data:   input.data,
            pos:    in.pos,
        }
        s, b, err := RunDFA(class, array, accept, fork)
        if b == nil {
            break
        }
        token := Token {}
        if k, ok := ac[s];ok {
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