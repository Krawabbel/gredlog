package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	input := flag.String("input", "/sys/class/thermal/thermal_zone0/temp", "path to the source")
	pattern := flag.String("pattern", "[0-9]*[\\.[0-9]*]?", "regular expression that matches the value to be logged")
	interval := flag.Duration("time", time.Second, "logging frequency")

	host := flag.String("host", "127.0.0.1", "host address of the REDIS database")
	port := flag.Int("port", 6379, "port of the REDIS database")
	key := flag.String("key", "redlog", "key for data in REDIS database")
	attempts := flag.Int("attempts", 1, "max. number or re-connect attempts for REDIS database")

	verbose := flag.Bool("verbose", true, "print logged values")

	if *key == "" || strings.Contains(*key, " ") {
		log.Fatal("key must not be empty or contain blanks")
	}
	flag.Parse()

	log.Println("Running redlog with arguments")
	flag.VisitAll(func(f *flag.Flag) { fmt.Printf("    %-8s: %v\n", f.Name, f.Value) })

	var src Source
	switch {
	case filepath.IsAbs(*input) || filepath.IsLocal(*input):
		src = FileSource{Path: *input}
	default:
		log.Fatalf("cannot deduce source type of '%s'", *input)
	}

	log.Fatal(Run(src, *host, *port, *key, *interval, *pattern, *verbose, *attempts))
}
