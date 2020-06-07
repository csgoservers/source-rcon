package main

import (
	"fmt"

	rcon "github.com/csgoservers/source-rcon/pkg/net"
)

func main() {
	opts := rcon.Options{
		Host:     "127.0.0.1",
		Password: "1234",
		Port:     27025,
	}
	conn := rcon.New(&opts)
	defer conn.Close()

	result, err := conn.ExecCommand("status") //cvarlist
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}
