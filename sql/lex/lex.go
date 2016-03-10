package lex

import (
    "../../log"
)

var (
    accept []bool
    array [][]int
    class []int
    ac map[string][]int
    dfa *preDfa
)

type Token struct {
    Kind string
    Raw []byte
}

func DefineToken(token string, g *nfa) {
    if _, ok := ac[token];ok {
        return
    }
    if len(ac) == 0 {
        fork := g.fork()
        fork.Simplify()
        dfa = fork.toDfa()
        ac[token] = make([]int, 0)
        for i, v := range dfa.v {
            if v {
                ac[token] = append(ac[token], i)
            }
        }
    } else {
        fork := g.fork()
        fork.Simplify()
        ac[token] = dfa.addDfa(fork.toDfa())
    }
}

func init() {
    ac = make(map[string][]int)
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

func defineTokens() {
    log.WriteLogSync("sys", "Start making DFA")
    DefineToken("float", links(repeat(numberNfa()), single('.'), chosable(repeat(numberNfa()))))
    DefineToken("int", repeat(numberNfa()))
    DefineToken("identical", links(repeat(letterNfa()), chosable(repeat(or(numberNfa(), letterNfa())))))
    DefineToken("split", or(single(' '), single('\t'), single('\n')))
    DefineToken("select", or(strings([]byte("select")), strings([]byte("SELECT"))))
    DefineToken("from", or(strings([]byte("from")), strings([]byte("FROM"))))
    DefineToken("where", or(strings([]byte("where")), strings([]byte("WHERE"))))
    DefineToken("and", or(strings([]byte("and")), strings([]byte("AND"))))
    DefineToken("or", or(strings([]byte("or")), strings([]byte("OR"))))
    DefineToken("order", or(strings([]byte("order")), strings([]byte("ORDER"))))
    DefineToken("by", or(strings([]byte("by")), strings([]byte("BY"))))
    DefineToken("insert", or(strings([]byte("insert")), strings([]byte("INSERT"))))
    DefineToken("update", or(strings([]byte("update")), strings([]byte("UPDATE"))))
    DefineToken("delete", or(strings([]byte("delete")), strings([]byte("DELETE"))))
    DefineToken("unReference", single('`'))
    DefineToken("reference", single('\''))
    DefineToken("dot", single(','))
    DefineToken("spot", single('.'))
    DefineToken("equal", single('='))
    DefineToken("lessEqual", link(single('<'), single('=')))
    DefineToken("moreEqual", link(single('>'), single('=')))
    DefineToken("more", single('>'))
    DefineToken("less", single('<'))
    DefineToken("notEqual", link(single('<'), single('>')))
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
        token := Token {}
        for k, v := range ac {
            for _, i := range v {
                if i == s {
                    token.Kind = k
                    goto escape
                }
            }
        }
        escape:
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