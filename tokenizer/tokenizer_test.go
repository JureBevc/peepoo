package tokenizer

import (
	"JureBevc/gpc/util"
	"fmt"
	"os"
	"testing"
)

func TestTokenizer(t *testing.T) {
	fmt.Println(os.Getwd())
	expected := []Token{
		"INTEGER",
		"PLUS",
		"INTEGER",
	}

	progList := []string{
		"../tests/p1.gpc",
		"../tests/p2.gpc",
		"../tests/p3.gpc",
		"../tests/p4.gpc",
	}

	tokenList := "../tests/tokens.list"

	for _, progPath := range progList {
		tokens := Tokenize(tokenList, progPath)
		if !util.CompareSlices(tokens, expected) {
			t.Errorf("Tokens not equal for tokens %s and program %s.\nActual:\n%s\nExpected:\n%s", tokenList, progPath, tokens, expected)
		}
	}

}
