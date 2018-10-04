package main

import (
	"flag"
)

func main() {

	// flag.Parse will modify values in Config.  Config and
	// default values can be seen in Config
	flag.Parse()

	srv := NewServer(":4444")

	done := make(chan bool)
	srv.Start(done)

	<-done
}
