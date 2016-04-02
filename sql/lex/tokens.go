package lex

import (
	"../../log"
)

func defineTokens() {
	log.WriteLogSync("sys", "Start making DFA")
	keyword := stringsToken("select", "from", "where", "update", "delete", "create", "insert", "into", "table", "order", "by")
	logical := stringsToken("and", "or", "not")
	structs := stringsToken("(", ")", ";", ",", ".")
	split := stringsToken(" ", "\t", "\n", "\r\n", "\r")
	relations := stringsToken(">", "<", ">=", "<=", "=", "<>")
	types := stringsToken("int", "float", "string", "object", "array")
	attributes := stringsToken("primary key", "unique")
	DefineCommon("floatval", links(repeat(numberNfa()), single('.'), chosable(repeat(numberNfa()))))
	DefineCommon("intval", repeat(numberNfa()))
	DefineCommon("identical", identicalNfa())
	DefineToken("keyword", keyword)
	DefineToken("types", types)
	DefineToken("logical", logical)
	DefineToken("structs", structs)
	DefineToken("split", split)
	DefineToken("relations", relations)
	DefineToken("unReference", single('`'))
	DefineToken("reference", single('\''))
	DefineToken("attributes", attributes)
	log.WriteLogSync("sys", "DFA prepared")
}
