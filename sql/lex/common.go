package lex

func digitNfa() *nfa {
	nfas := make([]*nfa, 0)
	for i := 0; i <= 9; i++ {
		nfas = append(nfas, single(byte(int('0')+i)))
	}
	return or(nfas...)
}

func integerNfa() *nfa {
	return repeat(digitNfa())
}

func floatNfa() *nfa {
	return links(integerNfa(), single('.'), integerNfa())
}

func letterNfa() *nfa {
	nfas := make([]*nfa, 0)
	for i := 0; i < 26; i++ {
		nfas = append(nfas, single(byte(int('a')+i)))
	}
	for i := 0; i < 26; i++ {
		nfas = append(nfas, single(byte(int('A')+i)))
	}
	return or(nfas...)
}

func identicalNfa() *nfa {
	return links(or(single('_'), letterNfa()), chosable(repeat(or(letterNfa(), digitNfa()))))
}
