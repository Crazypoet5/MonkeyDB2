package main

import (
    "../../lex"
    "fmt"
    "strings"
    "os"
)

func loop() bool {
    fmt.Print("Monkey>>")
    str := ""
    for !strings.Contains(str, ";") {
        buff := ""
        fmt.Scanf("%s", &buff)
        str += " " + buff
    }
    if strings.Contains(str, "quit") {
        return false
    }
    ts, err := lex.Parse(*lex.NewByteReader([]byte(str)))
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
    }
    fmt.Println(ts)
    return true
}

func main() {
    for loop() {}
}