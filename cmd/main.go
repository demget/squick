package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/demget/squick"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("squick: ")
	flag.Usage = func() { fmt.Fprintln(os.Stderr, help) }

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	sq, err := squick.New()
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-2]
	switch cmd {
	case "init":
		usage := func() { fmt.Fprintln(os.Stderr, helpInit) }
		if len(os.Args) <= 2 {
			usage()
			return
		}

		driver := os.Args[len(os.Args)-1]
		force := flag.Bool("force", false, "")
		pkg := flag.String("package", "database", "")
		flag.Parse()

		if driver == "" {
			usage()
			return
		}
		if *pkg == "" {
			log.Fatal("package option cannot be empty")
		}

		if *force {
			if err := os.RemoveAll(*pkg); err != nil {
				log.Fatal(err)
			}
		}

		if err := sq.Init(driver, *pkg); err != nil {
			log.Fatal(err)
		}
	case "make":
		usage := func() { fmt.Fprintln(os.Stderr, helpMake) }
		if len(os.Args) <= 2 {
			usage()
			return
		}

		stmt, err := squick.Parse(os.Args[len(os.Args)-1])
		if err != nil {
			log.Fatal(err)
		}

		name := flag.String("name", stmt.Model(), "")
		flag.Parse()

		if *name == "" {
			log.Fatal("name option cannot be empty")
		}

		if err := sq.Make(); err != nil {
			log.Fatal(err)
		}
	}
}
