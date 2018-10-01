package main

import (
	"flag"
	"fmt"
)

type Configuration struct {
	Addrport string

	Depth  int // walk depth
	Daemon bool
}

var (
	config Configuration
)

func init() {
	flag.StringVar(&config.Addrport, "addr", ":4444", "Address and port to run service on ")
	flag.BoolVar(&config.Daemon, "daemon", true, "Run as a deamon")
	flag.IntVar(&config.Depth, "depth", 1, "Walk Depth")
}

func main() {
	flag.Parse()

	// See what args are left
	if len(flag.Args()) > 0 {
		// additional args, we'll assume they are commands
	}

	fmt.Println("Starting inv server on ", config.Addrport)
	HTTPServer(config.Addrport)
}
