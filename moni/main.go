package main

import (
	"flag"
	"fmt"

	"github.com/rustyeddy/inv"
)

var (
	config inv.Configuration
)

func main() {
	config = inv.Config

	flag.Parse()

	var done chan bool
	switch {
	case config.Daemon:
		done = make(chan bool)
		inv.StartServer(config.Addrport, done)
	}

	_ = <-done
	fmt.Println("Server has finished")
}
