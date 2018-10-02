package main

import (
	"flag"
	"fmt"

	"github.com/rustyeddy/inv"
)

var (
	config *inv.Configuration
)

func main() {
	// Parse command line args setting config values
	// as set in config.go
	flag.Parse()

	// Copy of our config from inventory config
	config = inv.GetConfiguration()

	// Declare the done channel to communicate when the
	// server has completed
	var done chan bool

	// Figure what command we are going to run
	// with what specific arguments.
	switch {
	case config.Client:
		go inv.StartClient(config.Addrport, done)
	default:
		go inv.StartServer(config.Addrport, done)
	}

	_ = <-done
	fmt.Println("Server has finished")
}
