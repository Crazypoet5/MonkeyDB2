package lex

import "errors"

type statedDfa struct {
	class  []int
	array  [][]int
	accept []bool
	state  map[int]string
}

func (n *statedNfa) ToDFA() *statedDfa {
	n.Simplify()
	pre, dfa2nfa := n.g.toDfa()
	ret := &statedDfa{
		state: make(map[int]string),
	}
	for k, v := range dfa2nfa {
		if !pre.v[k] {
			continue
		}
		for _, i := range v {
			if !n.g.v[i] {
				continue
			}
			if kind, ok := ret.state[k]; !ok {
				ret.state[k] = n.state[i]
			} else if n.level[kind] < n.level[n.state[i]] {
				ret.state[k] = n.state[i]
			}
		}
	}
	ret.class, ret.array, ret.accept = pre.toArray()
	//	fmt.Println(ret.state)
	return ret
}

func (d *statedDfa) Parse(input ByteReader) ([]Token, error) {
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
		fork2 := ByteReader{
			data: input.data,
			pos:  in.pos,
		}
		tc, _ := fork2.Read()
		if tc == '\'' {
			b := make([]byte, 0)
			tc, _ = fork2.Read()
			for tc != '\'' || (len(b) > 0 && b[len(b)-1] == '\\') {
				if tc == '\'' {
					b[len(b)-1] = '\''
				} else {
					b = append(b, tc)
				}
				tc, _ = fork2.Read()
			}
			if tc == 0 {
				return nil, errors.New("No string end.")
			}
			token := Token{
				Kind: "string",
				Raw:  b,
			}
			t = append(t, token)
			in.pos = fork2.pos
			continue
		}
		s, b, err := RunDFA(d.class, d.array, d.accept, fork)
		if b == nil {
			break
		}
		token := Token{}
		if k, ok := d.state[s]; ok {
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
