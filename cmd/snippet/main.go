package main

import (
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/pkg/errors"
	"github.com/snamber/snippet"
	"log"
	"os"
	"reflect"
)

func main() {
	snippetsName := "snippets"
	templatesName := "templates"
	outputName := "output"
	var snippets string
	var templates string
	var output string

	flaggy.String(&snippets, "s", snippetsName, "The file or directory in which snippets are defined.")
	flaggy.String(&templates, "t", templatesName, "The file or directory containing the output template(s).")
	flaggy.String(&output, "o", outputName, "The directory where output file(s) will be written.")
	flaggy.Parse()

	fatalOnDefault(snippetsName, snippets)
	fatalOnDefault(templatesName, templates)
	fatalOnDefault(snippetsName, snippets)

	err := snippet.ExtractAndExpandSnippets(snippets, templates, output)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "could not extract and expand snippets"))
	}
}

func fatalOnDefault(name string, value interface{}) {
	if value == reflect.Zero(reflect.TypeOf(value)).Interface() {
		flaggy.ShowHelp(fmt.Sprintf("Please provide command line parameter '%s'\n", name))
		os.Exit(1)
	}
}
