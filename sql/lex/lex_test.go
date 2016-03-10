package lex

import "testing"
import "fmt"

func Test_DefineTokens(t *testing.T) {
    g := links(single('a'), single('b'), single('c'), chosable(repeat(or(single('a'), single('b'), single('c')))), 
    strings([]byte("cba")))
    DefineToken("Test", g)
    c, a, acs := dfa.toArray()
    s, _, _ := RunDFA(c, a, acs, ByteReader { data: []byte("abcbbcba") })
    flag := 0
    for _, i := range ac["Test"] {
        if i == s {
            flag = 1
            break
        }
    }
    if flag != 1 {
        t.Error("fail")
    }    
}

func TestParse(t *testing.T) {
    raw := []byte("select `name` from `table` where a and 1.4")
    fmt.Println(Parse(ByteReader{data: raw}))
}