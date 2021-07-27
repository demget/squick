package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/demget/squick"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("squick: ")

	flag := flag.NewFlagSet("squick", flag.ExitOnError)
	flag.Usage = func() { fmt.Fprintln(os.Stderr, help) }

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	driver, ok := os.LookupEnv("SQUICK_DRIVER")
	if !ok {
		log.Fatal("SQUICK_DRIVER environment key is unset")
	}

	sq, err := squick.New()
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, helpInit)
			return
		}

		force := flag.Bool("force", false, "")
		maxOpen := flag.Int("max-open", 0, "")
		maxIdle := flag.Int("max-idle", 2, "")
		flag.Parse(os.Args[2:])

		pkg := "database"
		if args := flag.Args(); len(args) > 0 {
			pkg = args[0]
		}

		if *force {
			if err := os.RemoveAll(pkg); err != nil {
				log.Fatal(err)
			}
		}

		if err := os.WriteFile(".squick", []byte(pkg), 0700); err != nil {
			log.Fatal(err)
		}

		ctx := squick.Context{
			MaxOpen: *maxOpen,
			MaxIdle: *maxIdle,
			Driver:  driver,
			Package: pkg,
		}
		if err := sq.Init(ctx); err != nil {
			log.Fatal(err)
		}
	case "make":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, helpMake)
			return
		}

		dburl, ok := os.LookupEnv("SQUICK_URL")
		if !ok {
			log.Fatal("SQUICK_URL environment key is unset")
		}

		db, err := sqlx.Open(driver, dburl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		initfile, _ := os.ReadFile(".squick")
		pkg := string(initfile)
		if pkg == "" {
			log.Fatal(`.squick file not found (use "squick init"`)
		}

		verbose := flag.Bool("v", false, "")
		ignore := flag.Bool("ignore", false, "")
		name := flag.String("name", "", "")
		tags := flag.String("tags", "json", "")
		updated := flag.String("updated", "", "")
		table := flag.String("table", "*", "")
		flag.Parse(os.Args[2:])

		stmt := squick.Parse(*table, flag.Args())
		if *name == "" {
			*name = stmt.Model()
		}

		ctx := squick.Context{
			Verbose:      *verbose,
			Ignore:       *ignore,
			DB:           db,
			Package:      pkg,
			Model:        *name,
			Tags:         strings.Split(*tags, ","),
			UpdatedField: *updated,
		}
		if err := sq.Make(ctx, stmt); err != nil {
			log.Fatal(err)
		}
	}
}
