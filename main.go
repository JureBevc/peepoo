package main

import (
	"JureBevc/peepoo/parser"
	"JureBevc/peepoo/runtime"
	"JureBevc/peepoo/tokenizer"
	"JureBevc/peepoo/util"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

//go:embed config/tokens.list
var tokenFile embed.FS

//go:embed config/grammar.list
var grammarFile embed.FS

func main() {
	// Parse flags
	verbose := flag.Int("verbose", 0, "Enable verbose mode")
	encodeString := flag.Bool("encode", false, "Encode string")
	decodeString := flag.Bool("decode", false, "Decode string")
	flag.Parse()
	util.LogLevel = *verbose

	// Get positional argument (first non-flag argument)
	var inputFile string = ""
	args := flag.Args()
	if len(args) > 0 {
		inputFile = args[0]
	}

	if *encodeString {
		if len(args) > 0 {
			inputString := args[0]
			encoded := runtime.EncodeString(inputString)
			fmt.Println(encoded)
		}
		return
	}

	if *decodeString {
		if len(args) > 0 {
			inputString := args[0]
			decoded := runtime.DecodeString(inputString)
			fmt.Println(decoded)
		}
		return
	}

	if inputFile == "" {
		log.Fatalln("Failed to open program file, no file provided.")
	}

	if _, err := os.Stat(inputFile); err != nil {
		log.Fatalf("Failed to open program file %s\n.", inputFile)
	}

	totalStart := time.Now()

	util.Log(1, "-Running tokenizer")
	start := time.Now()
	tokenDefinitions, tokens := tokenizer.Tokenize(tokenFile, inputFile)
	util.Log(1, fmt.Sprintf("Tokenizer finished: %s\n", time.Since(start)))

	util.Log(3, fmt.Sprintln(tokens))

	util.Log(1, fmt.Sprintln("-Running parser"))
	start = time.Now()
	ptree := parser.Parse(tokenDefinitions, tokens, grammarFile)
	util.Log(1, fmt.Sprintf("Parser finished: %s\n", time.Since(start)))

	util.Log(1, fmt.Sprintf("-Done: %s\n", time.Since(totalStart)))

	if util.LogLevel > 3 {
		parser.PrintTree(ptree, "")
	}
	runtime.RunTree(ptree)
}
