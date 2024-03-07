package main

import (
	"JureBevc/gpc/tokenizer"
	"flag"
	"fmt"
)

func main() {

	inputFile := flag.String("file", "", "Path to the input file")
	tokenDefinitions := flag.String("tokens", "", "Path to token definition file")
	flag.Parse()

	tokens := tokenizer.Tokenize(*tokenDefinitions, *inputFile)
	fmt.Println(tokens)
}
