package lex

import (
	"container/list"
	"fmt"
)

type statedNfa struct {
	g     *nfa
	state map[int]string
	level map[string]int
}

func NewBasic() *statedNfa {
	ret := &statedNfa{
		g:     identicalNfa(),
		state: make(map[int]string),
		level: make(map[string]int),
	}
	acID := len(ret.g.v) - 1
	ret.state[acID] = "identical"
	ret.level["identical"] = 0
	ret.appendToken("intval", 1, integerNfa())
	ret.appendToken("floatval", 0, floatNfa())
	keyword := stringsToken("select", "from", "where", "update", "delete", "create", "insert", "into", "table", "order", "by", "values", "dump", "set", "drop")
	logical := stringsToken("and", "or", "not")
	structs := stringsToken("(", ")", ";", ",", ".")
	split := stringsToken(" ", "\t", "\n", "\r\n", "\r")
	relations := stringsToken(">", "<", ">=", "<=", "=", "<>")
	types := stringsToken("int", "float", "string", "object", "array")
	attributes := stringsToken("primary key", "unique")
	wildcard := stringsToken("*", "%")
	ret.appendToken("keyword", 2, keyword)
	ret.appendToken("types", 4, types)
	ret.appendToken("logical", 4, logical)
	ret.appendToken("structs", 4, structs)
	ret.appendToken("split", 4, split)
	ret.appendToken("relations", 4, relations)
	ret.appendToken("unReference", 4, single('`'))
	ret.appendToken("reference", 4, single('\''))
	ret.appendToken("attributes", 4, attributes)
	ret.appendToken("wildcard", 4, wildcard)
	return ret
}

func NewTest() *statedNfa {
	ret := &statedNfa{
		g:     strings([]byte("abc")),
		state: make(map[int]string),
		level: make(map[string]int),
	}
	acID := len(ret.g.v) - 1
	ret.state[acID] = "abc"
	ret.level["abc"] = 0
	ret.appendToken("abd", 1, strings([]byte("abd")))
	return ret
}

func (n *statedNfa) Print() {
	n.g.Print()
	fmt.Println(n.state)
}

func (n *statedNfa) appendToken(kind string, level int, g *nfa) {
	n.level[kind] = level
	v := make([]bool, 1)
	v[0] = false
	v = append(v, n.g.v...)
	state := make(map[int]string)
	for k, v := range n.state {
		state[k+1] = v
	}
	for l := n.g.e.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		newP := path{
			x: p.x + 1,
			y: p.y + 1,
			c: p.c,
		}
		l.Value = newP
	}
	n.g.e.PushBack(path{
		x: 0,
		y: 1,
		c: 0,
	})
	base := len(v)
	v = append(v, g.v...)
	for l := g.e.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		n.g.e.PushBack(path{
			x: p.x + base,
			y: p.y + base,
			c: p.c,
		})
	}
	for i := base; i < len(v); i++ {
		if v[i] {
			state[i] = kind
		}
	}
	n.g.addEdge(0, base, 0)
	n.g.v = v
	n.state = state
}

func copyEdgeEx(g *statedNfa, v int) {
	bag := []int{}
	makeBag(g.g, v, &bag)
	//If there is a accept in one's bag, one should be accept
	for _, i := range bag {
		if g.g.v[i] {
			g.g.v[v] = true
			g.state[v] = g.state[i]
		}
	}
	for l := g.g.e.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		if inBag(p.x, bag) && !inBag(p.y, bag) {
			g.g.addEdge(v, p.y, p.c)
		}
	}
}

//Clear the null edge. After this step nfa would have more than one accept state
func (g *statedNfa) Simplify() {
	//Find all valid state
	valid := make([]bool, len(g.g.v))
	valid[0] = true
	for l := g.g.e.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		if p.c != 0 {
			valid[p.y] = true
			markAccept(g.g, p.y, &valid)
		}
	}
	//Add neccessary edges
	for i, v := range valid {
		if v {
			copyEdgeEx(g, i)
		}
	}
	//Delete useless edges
	for i := 0; i < len(g.g.v); i++ {
		if !valid[i] {
			g.g.v[i] = false //Only make it not acceptable rather than delete it
			if _, ok := g.state[i]; ok {
				delete(g.state, i)
			}
		}
	}
	eNew := list.List{}
	for l := g.g.e.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		if !(p.c == 0 || !valid[p.x] || !valid[p.y]) {
			eNew.PushBack(path{
				x: p.x,
				y: p.y,
				c: p.c,
			})
		}
	}
	//g.e = eNew
	//Delete alone edge
	eNew2 := list.List{}
	for l := eNew.Front(); l != nil; l = l.Next() {
		p := l.Value.(path)
		v := p.x
		if p.x == 0 {
			eNew2.PushBack(path{
				x: p.x,
				y: p.y,
				c: p.c,
			})
			continue
		}
		flag := 0
		for ll := eNew.Front(); ll != nil; ll = ll.Next() {
			pp := ll.Value.(path)
			if pp.y == v {
				flag = 1
				break
			}
		}
		if flag != 0 {
			eNew2.PushBack(path{
				x: p.x,
				y: p.y,
				c: p.c,
			})
		}
	}
	g.g.e = eNew2
}
