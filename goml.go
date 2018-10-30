package main

import (
	"go/printer"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/neelance/goml/parser"
)

func main() {
	log.SetFlags(0)

	for _, file := range os.Args[1:] {
		src, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		fset := token.NewFileSet()
		in, err := parser.ParseFile(fset, file, src)
		if err != nil {
			scanner.PrintError(os.Stderr, err)
			os.Exit(1)
		}

		out, err := os.Create(file[:len(file)-2])
		if err != nil {
			log.Fatal(err)
		}

		if err := printer.Fprint(out, fset, in); err != nil {
			log.Fatal(err)
		}
	}
}
