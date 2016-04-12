package lex

import (
    "container/list"
    "fmt"
)

type path struct {
    x, y int
    c byte
}

type nfa struct {
    v []bool        //accept
    e list.List
}

func (g nfa) fork() *nfa {
    r := &nfa {
        v:  make([]bool, len(g.v)),
        e:  list.List{},
    }
    copy(r.v, g.v)
    for e := g.e.Front();e != nil;e = e.Next() {
        p := e.Value.(path)
        r.addEdge(p.x, p.y, p.c)
    }
    return r
}

func (g *nfa) addEdge(x, y int, c byte) {
    g.e.PushBack(path {
        x:  x,
        y:  y,
        c:  c,
    })
}

func (g *nfa) addVertex(accept bool) int {
    g.v = append(g.v, accept)
    return len(g.v) - 1
}

func single(char byte) *nfa {
    r := &nfa {}
    r.addVertex(false)
    r.addVertex(true)
    r.addEdge(0, 1, char)
    return r
}

// Merge g1, g2, two nfa simply by collect all vertex and all path, making path correct in the new graph
func merge(g1, g2 *nfa) *nfa {
    r := g1.fork()
    base := len(g1.v)
    r.v = append(r.v, g2.v...)
    for l := g2.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        r.e.PushBack(path {
            x:  p.x + base,
            y:  p.y + base,
            c:  p.c,
        })
    }
    return r
}

func link(g1, g2 *nfa) *nfa {
    base := len(g1.v)
    r := merge(g1, g2)
    r.v[base - 1] = false           //Assert the end is accept
    r.addEdge(base - 1, base, 0)
    return r
}

func strings(str []byte) *nfa {
    if len(str) == 1 {
        return single(str[0])
    }
    return link(single(str[0]), strings(str[1:]))
}

func links(gs ...*nfa) *nfa {
    if len(gs) == 1 {
        return gs[0]
    }
    return link(gs[0], links(gs[1:]...))
}

func or(gs ...*nfa) *nfa {
    r := &nfa {}
    r.addVertex(false)
    end := []int {}                 //The set of accept states of gs
    for _, g := range gs {
        start := len(r.v)
        r = merge(r, g)
        r.addEdge(0, start, 0)      //Add r.start -> g.start
        end = append(end, len(r.v) - 1)
    }
    accept := len(r.v)
    r.addVertex(true)
    for _, i := range end {
        r.addEdge(i, accept, 0)
        r.v[i] = false
    }
    return r
}

func repeat(g *nfa) *nfa {
    base := len(g.v)
    g2 := g.fork()
    r := merge(g, g2)
    r.addEdge(base - 1, base, 0)
    r.v[base - 1] = false
    r.addEdge(len(r.v) - 1, base - 1, 0)
    r.v[len(r.v) - 1] = false
    newAc := r.addVertex(true)  //Make accept state in end
    r.addEdge(base - 1, newAc, 0)
    r.addEdge(newAc - 1, newAc, 0)
    return r
}

func chosable (g *nfa) *nfa {
    r := g.fork()
    r.addEdge(0, len(r.v) - 1, 0)
    return r
}

func (g *nfa) Test(input *ByteReader, current int, show bool) bool {
    if g.v[current] == true && input.Empty() {
        return true
    }
    inputRaw := input.Fork()
    next, _ := input.Read()
    for e := g.e.Front();e != nil;e = e.Next() {
        if e.Value == nil {
            continue
        }
        edge := e.Value.(path)
        if edge.x != current {
            continue
        }
        if edge.c == 0 {
            if show {
                fmt.Println(current, "->", edge.y)
            }
            inputNew := inputRaw.Fork()
            if g.Test(inputNew, edge.y, show) {
                return true
            }
        }
        if input.pos <= len(input.data) && edge.c == next {
            inputNew := input.Fork()
            if show {
                fmt.Println("Read: ", string(next))
            }
            if show {
                fmt.Println(current, "->", edge.y)
            }
            if g.Test(inputNew, edge.y, show) {
                return true
            }
        }
    }
    return false   
}

func (g *nfa) Print() {
    for i, v := range g.v {
        fmt.Println("V ", i, ":", v)
    }
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if p.c == 0 {
            fmt.Println("E ", p.x, "->", p.y)
        } else {
            fmt.Println("E ", p.x, "-(", string(p.c), ")>", p.y)
        }
    }
}

func markAccept(g *nfa, v int, valid *[]bool) {
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if p.x == v && p.c == 0 {
            if g.v[p.y] == true {
                (*valid)[p.y] = true
            }
            markAccept(g, p.y, valid)
        }
    }
}

func inBag(n int, bag []int) bool {
    for _, i := range bag {
        if i == n {
            return true
        }
    }
    return false
}

func makeBag(g *nfa, v int, bag *[]int) {
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if p.x == v && p.c == 0 {
            *bag = append(*bag, p.y)
            makeBag(g, p.y, bag)
        }
    }
}

func copyEdge(g *nfa, v int) {
    bag := []int {}
    makeBag(g, v, &bag)
    //If there is a accept in one's bag, one should be accept
    for _, i := range bag {
        if g.v[i] {
            g.v[v] = true
        }
    }
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if inBag(p.x, bag) && !inBag(p.y, bag) {
            g.addEdge(v, p.y, p.c)
        }
    }
}

//Clear the null edge. After this step nfa would have more than one accept state
func (g *nfa) Simplify() {
    //Find all valid state
    valid := make([]bool, len(g.v))
    valid[0] = true
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if p.c != 0 {
            valid[p.y] = true
            markAccept(g, p.y, &valid)
        }
    }
    //Add neccessary edges
    for i, v := range valid {
        if v {
            copyEdge(g, i)
        }
    }
    //Delete useless edges
    for i := 0;i < len(g.v);i++ {
        if !valid[i] {
            g.v[i] = false      //Only make it not acceptable rather than delete it
        }
    }
    eNew := list.List {}
    for l := g.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        if !(p.c == 0 || !valid[p.x] || !valid[p.y]) {
            eNew.PushBack(path {
                x:  p.x,
                y:  p.y,
                c:  p.c,
            })
        }
    }
    //g.e = eNew
    //Delete alone edge
    eNew2 := list.List{}
    for l := eNew.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        v := p.x
        if p.x == 0 {
            eNew2.PushBack(path {
                x:  p.x,
                y:  p.y,
                c:  p.c,
            })
            continue
        }
        flag := 0
        for ll := eNew.Front();ll != nil;ll = ll.Next() {
            pp := ll.Value.(path)
            if pp.y == v {
                flag = 1
                break
            }
        }
        if flag != 0 {
            eNew2.PushBack(path {
                x:  p.x,
                y:  p.y,
                c:  p.c,
            })
        }
    }
    g.e = eNew2
}