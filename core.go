package main

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

func Run(src Source, host string, port int, key string, interval time.Duration, pattern string, verbose bool) error {

	db, err := NewClient(host, port)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Println("connected to REDIS database")

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("error compiling regular expression: %s", err)
	}
	tic := time.NewTicker(interval)
	for {
		t := <-tic.C
		err := step(src, key, t, re, verbose, db)
		if err != nil {
			log.Println(err)
		}
	}
}

func step(src Source, key string, t time.Time, re *regexp.Regexp, verbose bool, db Client) error {
	timestamp := t.UnixMilli()
	raw, err := src.Read()
	if err != nil {
		return err
	}
	val := re.FindString(raw)
	if val == "" && verbose {
		log.Printf("warning: could not find regexp '%s' in got '%s'\n", re.String(), raw)
	}
	id, err := store(db, key, timestamp, string(val))
	if err != nil {
		return fmt.Errorf("error storing value: %s", err)
	}
	if verbose {
		log.Printf("[%s] data: %s, time: %d -> id: %s", key, string(val), timestamp, id)
	}
	return nil
}

var id_validator = regexp.MustCompile(`^"([0-9]*-[0-9]*)"$`)

func store(c Client, key string, timestamp int64, val string) (string, error) {
	q := fmt.Sprintf("XADD %s * time %d data %s", key, timestamp, val)
	r, err := c.Request(q)
	if err != nil {
		return "", err
	}
	m := id_validator.FindStringSubmatch(r)
	if m == nil {
		return "", fmt.Errorf("unexpected response: expected regexp '%s', got '%s', regexp result = %v", id_validator.String(), r, m)
	}
	return m[1], nil
}
