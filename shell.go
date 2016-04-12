package main

import (
	"bufio"
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
	inReader := bufio.NewReader(os.Stdin)
	i := 0
	for !strings.Contains(command, ";") {
		if i > 0 {
			fmt.Print("      ->")
		}
		cmd, _ := inReader.ReadString('\n')
		command += cmd
		i++
	}
	if strings.Contains(command, "quit;") {
		return false
	}
	//	fmt.Println(command)
	ts, _ := lex.Parse(*lex.NewByteReader([]byte(command)))
	fmt.Println(ts)
	stn, err := syntax.Parser(syntax.NewTokenReader(ts))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return true
	}
	stn.Print(1)
	r, re, err := plan.DirectPlan(stn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return true
	}
	r.Print()
	fmt.Println(re.AffectedRows, " row(s) affected in ", (float64)(re.UsedTime)/1000000000, "s.")
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
