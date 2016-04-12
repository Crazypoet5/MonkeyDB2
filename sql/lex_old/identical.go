package lex

import "fmt"

func allCharsExcept(exc byte) []byte {
    ac := make([]byte, 0)
    for l := byte('a');l < 'z';l++ {
        if l == exc {
            continue
        }
        ac = append(ac, l)
    }
    for l := byte('A');l < 'Z';l++ {
        if l == exc {
            continue
        }
        ac = append(ac, l)
    }
    for l := byte('0');l < '9';l++ {
        if l == exc {
            continue
        }
        ac = append(ac, l)
    }
    if exc != byte('_') {
        ac = append(ac, '_')
    }
    return ac
}

func identicalNfa() *nfa {
    return links(repeat(or(letterNfa(), single('_'))), chosable(repeat(or(numberNfa(), letterNfa(), single('_')))))
}

func generateIdenticalExceptString(exc []byte) *nfa {
    fmt.Println(len(exc))
    defer fmt.Println("return")
    if len(exc) == 1 {
        chars := allCharsExcept(exc[0])
        g := single(chars[0])
        for i := 1;i < len(chars);i++ {
            g = or(g, single(chars[i]))
        }
        return links(g, identicalNfa())
    }
    return or(generateIdenticalExceptString([]byte{exc[0]}), 
    links(single(exc[0]), generateIdenticalExceptString(exc[1:])))
}

func generateIdenticalExceptStrings(exc [][]byte) *nfa {
    g := generateIdenticalExceptString(exc[0])
    for i := 1;i < len(exc);i++ {
        g = or(g, generateIdenticalExceptString(exc[i]))
    }
    fmt.Println("ok")
    return g
}