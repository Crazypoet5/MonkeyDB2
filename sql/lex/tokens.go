package lex

func stringsToken(tokens ...string) *nfa {
	ret := make([]*nfa, 0)
	for _, v := range tokens {
		ret = append(ret, strings([]byte(v)))
	}
	return or(ret...)
}

func defineTokens() {
	//	log.WriteLogSync("sys", "Start making DFA")
	//	keyword := stringsToken("select", "from", "where", "update", "delete", "create", "insert", "into", "table", "order", "by", "values", "dump")
	//	logical := stringsToken("and", "or", "not")
	//	structs := stringsToken("(", ")", ";", ",", ".")
	//	split := stringsToken(" ", "\t", "\n", "\r\n", "\r")
	//	relations := stringsToken(">", "<", ">=", "<=", "=", "<>")
	//	types := stringsToken("int", "float", "string", "object", "array")
	//	attributes := stringsToken("primary key", "unique")
	//	NFA.appendToken("keyword", 2, keyword)
	//	NFA.appendToken("types", 4, types)
	//	NFA.appendToken("logical", 4, logical)
	//	NFA.appendToken("structs", 4, structs)
	//	NFA.appendToken("split", 4, split)
	//	NFA.appendToken("relations", 4, relations)
	//	NFA.appendToken("unReference", 4, single('`'))
	//	NFA.appendToken("reference", 4, single('\''))
	//	NFA.appendToken("attributes", 4, attributes)
	//	log.WriteLogSync("sys", "DFA prepared")
}
