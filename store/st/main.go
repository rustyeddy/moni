package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rustyeddy/store"
)

var (
	basedir string
)

func init() {
	flag.StringVar(&basedir, "dir", ".", "Use path as a Store, cwd by default")
}

func main() {
	flag.Parse()

	st, err := store.UseStore(basedir)
	if err != nil {
		fmt.Println("Unable to use %s as storage")
		os.Exit(4)
	}

	// Create a writer to pass into our command handler
	wr := os.Stdout
	cmd := flag.Arg(0)

	switch cmd {
	case "ls":
		storeList(wr, st)
	case "cp":
		fmt.Println("Todo: cp <sobj> <dobj>")
	case "du":
		fmt.Println("Todo: du")
	case "cat":
		fmt.Println("Todo: cat <obj>")
	default:
		fmt.Println("unsupported command ", cmd)
	}
}

func storeList(w io.Writer, st *store.Store) {
	idx := st.GetIndex()
	if idx == nil {
		fmt.Println("failed to get index from store ", st.Path)
	}
	str := ""
	for n, o := range idx {
		str = n + " " + o.String() + "\n"
	}
	fmt.Fprintln(w, str)
}
