// Package waitforpg waits for Postgres to be accessible before quitting.
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx"
)

func usage() {
	log.Println("usage: waitforpg [-i 1s] [-t 10s] <database_url>")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	log.SetFlags(0)
	var (
		timeout  = flag.Duration("t", 3*time.Second, "timeout")
		interval = flag.Duration("i", 100*time.Millisecond, "interval between checks")
		verbose  = flag.Bool("v", false, "verbose mode")
	)
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	cfg, err := pgx.ParseConnectionString(flag.Arg(0))
	if err != nil {
		log.Fatalf("error parsing connection string: %v\n", err)
	}

	start := time.Now()
	for time.Since(start) < *timeout {
		_, err = pgx.Connect(cfg)
		if err == nil {
			if *verbose {
				log.Println("success.")
			}
			os.Exit(0)
		}
		if *verbose {
			log.Printf("err: %v\n", err)
		}
		time.Sleep(*interval)
	}
	log.Fatalf("could not connect to postgres: %v\n", err)
}
