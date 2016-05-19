package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"./network"
)

var port = flag.String("port", "2016", "Dial port")

func loop(tcpSession *network.TCPSession) bool {
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
	tcpSession.SendPack(network.MakePack(network.DIRECT_QUERY, []byte(command)))
	pack := tcpSession.RecvPack()
	switch pack.Type {
	case network.RESPONSE:
		fmt.Println(string(pack.Data))
	}
	return true
}

func main() {
	flag.Parse()
	tcpSession := network.Dial("127.0.0.1:" + *port)
	for loop(tcpSession) {

	}
}
