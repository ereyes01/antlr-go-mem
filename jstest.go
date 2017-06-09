package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"./parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var filename string

func init() {
	flag.StringVar(&filename, "file", "", "file to parse")
}

type JSScanner struct {
	*parser.BaseECMAScriptListener
}

func (s *JSScanner) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println("RULE:", ctx.GetText())
}

func main() {
	flag.Parse()
	if filename == "" {
		log.Fatalln("you must supply a filename")
	}

	// for the profiler
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	input, err := antlr.NewFileStream(filename)
	if err != nil {
		log.Fatal(err)
	}
	lexer := parser.NewECMAScriptLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewECMAScriptParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Program()
	antlr.ParseTreeWalkerDefault.Walk(new(JSScanner), tree)
}
