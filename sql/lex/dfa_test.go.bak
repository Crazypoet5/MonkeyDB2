package lex

import "testing"
import _"fmt"

func Test_toDfa(t *testing.T) {/*
    g := links(single('a'), single('b'), single('c'), chosable(repeat(or(single('a'), single('b'), single('c')))), 
    strings([]byte("cba")))
    g.Simplify()
    g.Print()
    g.toDfa()*/
}

func TestRunDFA(t *testing.T) {
    g := links(single('a'), single('b'), single('c'), chosable(repeat(or(single('a'), single('b'), single('c')))), 
    strings([]byte("cba")))
    g.Simplify()
    dfa := g.toDfa()
    //dfa.Print()
    c, a, ac := dfa.toArray()
    _, _, err := RunDFA(c, a, ac, ByteReader {
        data:   []byte("abcbsa"),
    })
    if err == nil {
        t.Error("abcbsa, s not checked")
    }
    _, _, err = RunDFA(c, a, ac, ByteReader {
        data:   []byte("abcbbbcba"),
    })
    if err != nil {
        t.Error("abcbbbcba with an error:" + err.Error())
    }
    _, _, err = RunDFA(c, a, ac, ByteReader {
        data:   []byte("abcbbbba"),
    })
    if err == nil {
        t.Error("abcbbbba success")
    }
    _, get, err := RunDFA(c, a, ac, ByteReader {
        data:   []byte("abcbbbcbaaa"),
    })
    if err == nil {
        t.Error("abcbbbcbaaa success")
    }
    if string(get) != "abcbbbcba" {
        t.Error("abcbbbcbaaa, abcbbbcba not get")
    }
    _, _, err = RunDFA(c, a, ac, ByteReader {
        data:   []byte("abcbbbcba"),
    })
    if err != nil {
        t.Error("abcbbbcba fail")
    }
}

func Test_findBag(t *testing.T) {
    g := links(strings([]byte("sel")), or(single('e'), single('p')), strings([]byte("ekk")))
    g.Simplify()
    dfa1 := g.toDfa()
    g2 := strings([]byte("seleee"))
    g2.Simplify()
    dfa2 := g2.toDfa()
    b1, b2 := findBag(dfa1, dfa2)
    if b1[3] != 3 || b1[5] != 6 || b2[5] != 5 {
        t.Error("fing bag error")
    }
}

func Test_addDfa(t *testing.T) {
    g := links(strings([]byte("sel")), or(single('e'), single('p')), strings([]byte("ekk")))
    g.Simplify()
    dfa1 := g.toDfa()
    g2 := strings([]byte("seleee"))
    g2.Simplify()
    dfa2 := g2.toDfa()
    if dfa1.addDfa(dfa2)[0] != 14 {
        t.Error("accept error")
    }/*
    f := links(repeat(numberNfa()), single('.'), chosable(repeat(numberNfa())))
    i := repeat(or(numberNfa()))
    f.Simplify()
    i.Simplify()
    dfa1 = f.toDfa()
    dfa2 = i.toDfa()
    dfa1.addDfa(dfa2)
    a, b, c := dfa1.toArray()
    fmt.Println(RunDFA(a, b, c, ByteReader{data:[]byte("1234566")}))
    fmt.Println(" --- ")*/
}