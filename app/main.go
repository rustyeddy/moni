package main

import (
	"flag"
	"fmt"
)

type Configuration struct {
	Addrport string
	Depth    int // walk depth
}

var (
	config Configuration
)

func init() {
	flag.StringVar(&config.Addrport, "addr", ":4444", "Address and port to run service on ")
}

func main() {
	flag.Parse()

	fmt.Println("Starting inv server on ", config.Addrport)
	HTTPServer(config.Addrport)
}
