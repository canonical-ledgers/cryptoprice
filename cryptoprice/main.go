package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/canonical-ledgers/cryptoprice"
)

func main() {
	flag.Parse()
	// Attempt to run the completion program.
	if completion.Complete() {
		// The completion program ran, so just return.
		return
	}
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	dt := time.Now()
	if _, ok := flagset["datetime"]; ok {
		dt = time.Time(datetime)
	}

	if debug {
		fmt.Println("from    ", fromSymbol)
		fmt.Println("to      ", toSymbol)
		fmt.Println("datetime", dt)
	}

	c := cryptoprice.NewClient(strings.ToUpper(fromSymbol), strings.ToUpper(toSymbol))
	c.Timeout = timeout
	p, err := c.GetPrice(dt)
	if err != nil {
		log.Fatalf("error: %#v\n", err)
	}
	fmt.Println(p)
}
