package main

import (
	"flag"
	"fmt"

	"./recovery"

	"./network"
)

var port = flag.String("port", "2016", "Listen port.")

func bye() {
	recovery.SafeExit()
	fmt.Println("Bye!")
}

func main() {
	flag.Parse()
	network.Listen(*port)
	bye()
}
