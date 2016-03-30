package lex

import (
    "container/list"
    "fmt"
    "errors"
)

type queue []int

type preDfa struct {
    v   []bool
    e   list.List
}

func (g *preDfa) addEdge(x, y int, c byte) {
    g.e.PushBack(path {
        x:  x,
        y:  y,
        c:  c,
    })
}

func (g *preDfa) addVertex(accept bool) int {
    g.v = append(g.v, accept)
    return len(g.v) - 1
}

func (g *preDfa) Print() {
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

func (g preDfa) fork() *preDfa {
    r := &preDfa {
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

func (q *queue) push(i int) {
    *q = append((*q), i)
}

func (q *queue) front() int {
    return (*q)[0]
}

func (q *queue) pop() {
    *q = (*q)[1:]
}

func (q *queue) empty() bool {
    return len(*q) == 0
}

func isSubSet(s []int, d []int) bool {
    if len(s) > len(d) {
        return false
    }
    for _, si := range s {
        flag := 0
        for _, di := range d {
            if si == di {
                flag = 1
                break
            }
        }
        if flag == 0 {
            return false
        }
    }
    return true
}

func (g *nfa) toDfa() *preDfa {
    dfa2nfa := make(map[int][]int)
    dfa := &preDfa {}
    L := queue{}
    L.push(0)
    ind := dfa.addVertex(false)
    dfa2nfa[ind] = []int{0}
    for !L.empty() {
        v := L.front()
        vs := dfa2nfa[v]
        L.pop()
        //计算从这个状态输出的所有边所接受的字符集的并集
        charset := []byte{}
        for l := g.e.Front();l != nil;l = l.Next() {
            p := l.Value.(path)
            if isSubSet([]int{p.x}, vs) {
                flag := 0
                for _, c := range charset {
                    if c == p.c {
                        flag = 1
                        break
                    }
                }
                if flag == 0 {
                    charset = append(charset, p.c)
                }
            }
        }
        //然后对该集合中的每一个字符寻找接受这个字符的边，把这些边的目标状态的并集T计算出来
        for _, c := range charset {
            T := []int{}
            for l := g.e.Front();l != nil;l = l.Next() {
                p := l.Value.(path)
                if p.c == c && isSubSet([]int{p.x}, vs) {
                    T = append(T, p.y)
                }
            }
            //如果T∈D则代表当前字符指向了一个已知的DFA状态。否则 则代表当前字符指向了一个未创建的DFA状态，这个时候就把T放进L和D中。
            flag := 0
            for i, m := range dfa2nfa {
                if isSubSet(T, m) && isSubSet(m, T) {
                    flag = 1
                    dfa.addEdge(v, i, c)
                    break
                }
            }
            if flag == 0 {
                //fmt.Println("New state:", T)
                b := false
                for _, t := range T {
                    if g.v[t] {
                        b = true
                    }
                }
                i := dfa.addVertex(b)
                dfa.addEdge(v, i, c)
                dfa2nfa[i] = T
                L.push(i)
            }
        }
    }  
    return dfa
}

//Class [ASCII] = CLASS
//Array [CurrentState][Char] = NextState
//Accept [CurrentState] = Bool
func (dfa *preDfa) toArray() ([]int, [][]int, []bool) {
    charset := make([]byte, 0)
    class := make([]int, 256)
    for i := 0;i < 256;i++ {
        class[i] = -1
    }
    for l := dfa.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        flag := 0
        for _, c := range charset {
            if c == p.c {
                flag = 1
                break
            }
        }
        if flag == 0 {
            charset = append(charset, p.c)
            class[p.c] = len(charset) - 1
        }
    }
    array := make([][]int, len(dfa.v))
    for i := 0;i < len(dfa.v);i++ {
        array[i] = make([]int, len(charset))
        for j := 0;j < len(charset);j++ {
            array[i][j] = -1
        }
    }
    for l := dfa.e.Front();l != nil;l = l.Next() {
        p := l.Value.(path)
        array[p.x][class[p.c]] = p.y
    }
    accept := make([]bool, len(dfa.v))
    for i, v := range dfa.v {
        if v {
            accept[i] = true
        }
    }
    return class, array, accept
}


func RunDFA(class []int, array [][]int, accept []bool, reader ByteReader) (int, []byte, error) {
    last := -1
    current := 0
    get := make([]byte, 0)
    accepted := false
    lastAccepted := false
    lastGet := make([]byte, 0)
    for {
        c, err := reader.Read()
        if err != nil {
            if accepted {
                return current, get, nil
            } else {
                if lastAccepted {
                    return current, lastGet, errors.New("Unexpected char:" + string(get[len(lastGet)]))
                }
                return current, nil, errors.New("Unexpected char:" + string(get[0]))
            }
        }
        if class[c] == -1 {
            if lastAccepted {
                return current, lastGet, errors.New("Unexpected char:" + string(c))
            }
            return current, nil, errors.New("Unexpected char:" + string(c))
        }
        last = current
        current = array[current][class[c]]
        if current == -1 {
            if accepted {
                return last, get, errors.New("Unexpected char:" + string(c))
            } else {
                if lastAccepted {
                    return last, lastGet, errors.New("Unexpected char:" + string(get[len(lastGet)]))
                }
                return last, nil, errors.New("Unexpected char:" + string(get[0]))
            }
        }
        if accept[current] {
            get = append(get, c)
            accepted = true
            lastGet = []byte{}
            lastGet = append(lastGet, get...)
            lastAccepted = true
        } else {
            get = append(get, c)
            accepted = false
        }
    }
}

//find bags for each with each own id
func findBag(g1, g2 *preDfa) ([]int, []int) {
    q1, q2 := queue{}, queue{}
    q1.push(0)
    q2.push(0)
    set1, set2 := []int{0}, []int{0}
    for !q1.empty() && !q2.empty() {
        v1 := q1.front()
        q1.pop()
        v2 := q2.front()
        q2.pop()
        for l1 := g1.e.Front();l1 != nil;l1 = l1.Next() {
            p1 := l1.Value.(path)
            if p1.x != v1 {
                continue
            }
            for l2 := g2.e.Front();l2 != nil;l2 = l2.Next() {
                p2 := l2.Value.(path)
                if p2.x != v2 {
                    continue
                }
                if p1.c == p2.c {
                    if !_inBag(p1.y, set1) && !_inBag(p2.y, set2) {
                        q1.push(p1.y)
                        q2.push(p2.y)
                        set1 = append(set1, p1.y)
                        set2 = append(set2, p2.y)
                    }
                    
                }
            }
        }
    }
    return set1, set2
}

func _inBag(v int, bag []int) bool {
    for _, i := range bag {
        if i == v {
            return true
        }
    }
    return false
} 

func convert(bag1, bag2 []int) []int {
    bag1max := bag1[0]
    for i := range bag1 {
        if bag1[i] > bag1max {
            bag1max = bag1[i]
        }
    }
    r := make([]int, bag1max + 1)
    for i := range bag1 {
        r[bag1[i]] = bag2[i]
    }
    return r
}

//add additional dfa to a exist dfa, and return accept states of additional dfa
func (g *preDfa) addDfa(add *preDfa) []int {
    bag1, bag2 := findBag(g, add)
    convert := convert(bag2, bag1)
    base := len(g.v) - 1
    g.v = append(g.v, add.v[1:]...)
    accept := make([]int, 0)
    for i, v := range add.v {
        if v {
            accept = append(accept, i + base)
        }
    }
    for i, _ := range bag1 {
        if add.v[bag2[i]] {
            g.v[bag1[i]] = true
            accept = append(accept, bag1[i])
        }
    }
    
    eNew := list.List{}
    for al := add.e.Front();al != nil;al = al.Next() {
        ap := al.Value.(path)
        xInB, yInB := _inBag(ap.x, bag2), _inBag(ap.y, bag2)
        if xInB && yInB {
            continue
        }
        if xInB {
            eNew.PushBack(path {
                x:  convert[ap.x],
                y:  base + ap.y,
                c:  ap.c,
            })
        } else if yInB {
            eNew.PushBack(path {
                x:  base + ap.x,
                y:  convert[ap.y],
                c:  ap.c,
            })
        } else {
            eNew.PushBack(path {
                x:  base + ap.x,
                y:  base + ap.y,
                c:  ap.c,
            })
        }
    }
    g.e.PushBackList(&eNew)
    
    return accept
}
