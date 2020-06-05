package main

import (
	"log"

	rcon "github.com/csgoservers/source-rcon/pkg/net"
)

func main() {
	opts := rcon.Options{
		Host: "yybg.counter.monster",
		Port: 27015,
	}
	conn := rcon.New(&opts)
	defer conn.Close()

	_, err := conn.ExecCommand("stats")
	log.Println(err)
}
