package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mehiX/cli-xkcd/client"
)

var fn = flag.String("in", "", "JSON input file")

/**
XKCD data is cached locally in JSON
This tool shows the contents of the cache.
Receives input from standard input or reads a file passed in as parameter
**/
func main() {

	flag.Parse()

	var data []byte
	var err error

	if *fn != "" {
		data, err = ioutil.ReadFile(*fn)
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		panic(err)
	}

	results := make([]*client.Result, 0)

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&results); err != nil {
		panic(err)
	}

	prettyPrint(results)
}

func prettyPrint(res []*client.Result) {

	template.Must(template.New("results").Parse(tmpl)).Execute(os.Stdout, struct {
		TotalCount int
		Results    []*client.Result
	}{
		TotalCount: len(res),
		Results:    res,
	})
}

const tmpl = `{{.TotalCount}} issues:
{{range .Results}}-------------------------------
Number: {{.Num}}
Date: {{.Day}}/{{.Month}}/{{.Year}}
Title: {{.Title}}
{{end}}
`
