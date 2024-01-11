package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	src := flag.String("src", "/sys/class/thermal/thermal_zone0/temp", "path to the source")
	pattern := flag.String("pattern", "[0-9]*[\\.[0-9]*]?", "regular expression that matches the value to be logged")
	interval := flag.Duration("time", time.Second, "logging frequency")

	host := flag.String("host", "127.0.0.1", "host address of the REDIS database")
	port := flag.Int("port", 6379, "port of the REDIS database")
	key := flag.String("key", "gredlog", "key for data in REDIS database")

	verbose := flag.Bool("verbose", true, "print logged values")

	flag.Parse()

	log.Println("Running gredlog with arguments")
	flag.VisitAll(func(f *flag.Flag) { fmt.Printf("    %-8s: %v\n", f.Name, f.Value) })

	log.Fatal(Run(*src, *host, *port, *key, *interval, *pattern, *verbose))
}
