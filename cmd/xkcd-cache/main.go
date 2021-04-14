package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mehiX/cli-xkcd/client"
)

var workers = flag.Int("w", 60, "Number of concurrent fetch routines")
var outFile = flag.String("out", "", "Output JSON results to this file, otherwise standard output")

func main() {

	flag.Parse()

	if *workers <= 0 {
		fmt.Println("Invalid number of workers")
		flag.Usage()
	}

	out := os.Stdout
	if *outFile != "" {
		// try to use a file for output
		var err error
		if out, err = os.OpenFile(*outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
			panic(err)
		}
	}

	client.FetchAll(out, *workers)
}
