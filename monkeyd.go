package main

import (
	"fmt"

	"./recovery"

	"./network"
)

func bye() {
	recovery.SafeExit()
	fmt.Println("Bye!")
}

func main() {
	network.Listen()
	bye()
}
