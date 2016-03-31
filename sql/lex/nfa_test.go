package lex

import "testing"

func Test_merge(t *testing.T) {
    g := single('a')
    g2 := single('b')
    g3 := merge(g, g2)
    if len(g3.v) != 4 {
        t.Error("merge len v error")
    }
    if g3.e.Len() != 2 {
        t.Error("merge len e error")
    }
}

func Test_link(t *testing.T) {
    g := single('a')
    g2 := single('b')
    g3 := link(g, g2)
    if len(g3.v) != 4 {
        t.Error("link len v error")
    }
    if g3.e.Len() != 3 {
        t.Error("link len e error")
    }
}

func Test_or(t *testing.T) {
    g := single('a')
    g2 := single('b')
    g4 := single('c')
    g3 := link(g, g2)
    g5 := link(g, g4)
    g6 := or(g3, g5)
    if len(g6.v) != 10 {
        t.Error("or len v error")
    }
}

func Test_repeat(t *testing.T) {
    g := single('a')
    g2 := single('c')
    g3 := link(g, g2)
    g4 := repeat(g3)
    if len(g4.v) != 9 {
        t.Error("repeat len v error")
    }
}

func Test_chosable(t *testing.T) {
    g := single('a')
    g2 := single('c')
    g3 := chosable(g2)
    g4 := link(g, g3)
    if len(g4.v) != 4 {
        t.Error("or len v error")
    }
}

func TestTest(t *testing.T) {
    g := links(single('a'), single('b'), single('c'), chosable(repeat(or(single('a'), single('b'), single('c')))), 
    strings([]byte("cba")))
    input := &ByteReader {
        data:   []byte("abcabcccccba"),
    }
    if !g.Test(input, 0, false) {
        t.Error("abcabcccccba fail")
    }
    input = &ByteReader {
        data:   []byte("abccba"),
    }
    if !g.Test(input, 0, false) {
        t.Error("abccba fail")
    }
    input = &ByteReader {
        data:   []byte("abcba"),
    }
    if g.Test(input, 0, false) {
        t.Error("abcba success")
    }
    input = &ByteReader {
        data:   []byte("acba"),
    }
    if g.Test(input, 0, false) {
        t.Error("acba success")
    }
}

func TestSimplify(t *testing.T) {
    g := links(single('a'), or(links(single('c'), chosable(links(single('c'), single('e'), single('p'), single('t')))), 
    links(single('b'), single('o'), single('r'), single('t'))))
    g.Simplify()
    input := &ByteReader {
        data:   []byte{'a', 'c', 'c'},
    }
    if g.Test(input, 0, false) {
        t.Error("acc success")
    }
    input = &ByteReader {
        data:   []byte{'a', 'c'},
    }
    if !g.Test(input, 0, false) {
        t.Error("ac fail")
    }
    input = &ByteReader {
        data:   []byte{'a', 'c', 'c', 'e', 'p', 't'},
    }
    if !g.Test(input, 0, false) {
        t.Error("accept fail")
    }
    input = &ByteReader {
        data:   []byte{'a', 'b', 'c', 'e', 'p', 't'},
    }
    if g.Test(input, 0, false) {
        t.Error("abcept success")
    }
    input = &ByteReader {
        data:   []byte{'a', 'b', 'o', 'r', 't'},
    }
    if !g.Test(input, 0, false) {
        t.Error("abort fail")
    }
    
    g = links(single('a'), single('b'), single('c'), chosable(repeat(or(single('a'), single('b'), single('c')))), 
    strings([]byte("cba")))
    g.Simplify()
    input = &ByteReader {
        data:   []byte("abcabcabcabccba"),
    }
    if !g.Test(input, 0, false) {
        t.Error("abcabcabcabccba fail")
    }
    input = &ByteReader {
        data:   []byte("abccba"),
    }
    if !g.Test(input, 0, false) {
        t.Error("abccba fail")
    }
    input = &ByteReader {
        data:   []byte("abcba"),
    }
    if g.Test(input, 0, false) {
        t.Error("abcba success")
    }
    input = &ByteReader {
        data:   []byte("acba"),
    }
    if g.Test(input, 0, false) {
        t.Error("acba success")
    }
}