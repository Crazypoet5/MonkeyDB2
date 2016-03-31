package main

import (
	"fmt"
	"os"
	"strings"

	"./plan"
	"./sql/lex"
	"./sql/syntax"
)

func welcome() {
	fmt.Println("MonkeyDB2 @ 2016")
	fmt.Println("Welcome to SQL play ground~")
}

func loop() bool {
	fmt.Print("Monkey>>")
	command := ""
	for !strings.Contains(command, ";") {
		buff := ""
		fmt.Scanf("%s", &buff)
		command += buff + " "
	}
	if command == " quit; " {
		return false
	}
	//fmt.Println(command)
	ts, _ := lex.Parse(*lex.NewByteReader([]byte(command)))
	//fmt.Println(ts)
	stn, err := syntax.Parser(syntax.NewTokenReader(ts))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		stn.Print(1)
	}
	fmt.Println(plan.CreatePlan(stn))
	return true
}

func bye() {
	fmt.Println("Bye!")
}

func main() {
	welcome()
	for loop() {
	}
	bye()
}
