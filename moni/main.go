package main

import (
	"flag"
	"fmt"

	"github.com/rustyeddy/inventory/inv"
)

var (
	config inv.Configuration
)

func main() {
	config = inv.Config

	flag.Parse()

	fmt.Printf("Starting it up %+v\n", config)
	var done chan bool
	switch {
	case config.Daemon:
		done = make(chan bool)

		fmt.Println("Running new server...")
		server := inv.NewServer(config.Addrport)
		go server.Start(done)
	}

	_ = <-done
	fmt.Println("Server has finished")
}
