package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() { fmt.Fprintln(os.Stderr, help) }
	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, helpInit)
			return
		}
	case "make":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, helpMake)
			return
		}
	}
}