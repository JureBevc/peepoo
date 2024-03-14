package main

import (
	"JureBevc/gpc/parser"
	"JureBevc/gpc/tokenizer"
	"flag"
	"fmt"
)

func main() {

	inputFile := flag.String("file", "prog.gpc", "Path to the input file")
	tokensFile := flag.String("tokens", "tokens.list", "Path to token definition file")
	grammarFile := flag.String("grammar", "grammar.list", "Path to grammar definition file")
	flag.Parse()

	fmt.Printf("#Loading configuration:\nProgram=%s\nTokens=%s\nGrammar=%s\n", *inputFile, *tokensFile, *grammarFile)

	tokenDefinitions, tokens := tokenizer.Tokenize(*tokensFile, *inputFile)

	fmt.Println("#Tokenizer result:")
	fmt.Println(*tokens)

	grammar := parser.Parse(tokenDefinitions, tokens, *grammarFile)

	fmt.Println("#Parser results:")
	fmt.Println(*grammar)
}
