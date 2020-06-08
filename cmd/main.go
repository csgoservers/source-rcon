package main

import (
	"flag"
	"fmt"

	rcon "github.com/csgoservers/source-rcon/pkg/protocol"
)

var (
	host     = flag.String("H", "127.0.0.1", "Host where game server is running")
	port     = flag.Int("p", 27015, "Port where server is listening for connection")
	password = flag.String("s", "", "RCON secret password")
	cmd      = flag.String("c", "", "Command to send to server")
)

func main() {
	flag.Parse()

	if *cmd == "" {
		fmt.Println("You need to specify a command.")
		return
	}
	opts := rcon.Options{
		Host:     *host,
		Password: *password,
		Port:     *port,
	}
	conn := rcon.New(&opts)
	defer conn.Close()

	result, err := conn.ExecCommand(*cmd) //cvarlist
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
