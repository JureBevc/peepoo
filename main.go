package main

import (
	"JureBevc/gpc/assembler"
	"JureBevc/gpc/parser"
	"JureBevc/gpc/tokenizer"
	"flag"
	"fmt"
	"path/filepath"
	"time"
)

func main() {
	rootConfigPath := flag.String("config", ".", "Path to root of config files")
	inputFile := flag.String("file", "prog.gpc", "Path to the input file")
	tokensFile := flag.String("tokens", "tokens.list", "Path to token definition file")
	grammarFile := flag.String("grammar", "grammar.list", "Path to grammar definition file")
	assemblerFile := flag.String("assemble", "assemble.list", "Path to assembler definition file")
	flag.Parse()

	*inputFile = filepath.Join(*rootConfigPath, *inputFile)
	*tokensFile = filepath.Join(*rootConfigPath, *tokensFile)
	*grammarFile = filepath.Join(*rootConfigPath, *grammarFile)
	*assemblerFile = filepath.Join(*rootConfigPath, *assemblerFile)

	totalStart := time.Now()

	fmt.Println("-Running tokenizer")
	start := time.Now()
	tokenDefinitions, tokens := tokenizer.Tokenize(*tokensFile, *inputFile)
	fmt.Printf("Tokenizer finished: %s\n", time.Since(start))

	fmt.Println("-Running parser")
	start = time.Now()
	ptree := parser.Parse(tokenDefinitions, tokens, *grammarFile)
	fmt.Printf("Parser finished: %s\n", time.Since(start))

	fmt.Println("-Running assembler")
	start = time.Now()
	assembler.Assemble(ptree, *assemblerFile)
	fmt.Printf("Assembler finished: %s\n", time.Since(start))

	fmt.Printf("-Done: %s\n", time.Since(totalStart))
}
