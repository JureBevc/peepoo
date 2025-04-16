package tokenizer

import (
	"fmt"
	"os"
	"testing"
)

func compareTokens(t1 *[]Token, t2 *[]Token) bool {
	if len(*t1) != len(*t2) {
		return false
	}

	for i := 0; i < len(*t1); i++ {
		if (*t1)[i].Name != (*t2)[i].Name {
			return false
		}
	}

	return true
}

func TestTokenizer(t *testing.T) {
	fmt.Println(os.Getwd())
	expected := []Token{
		{Name: "integer"},
		{Name: "plus"},
		{Name: "integer"},
		{Name: "seperator"},
	}

	progList := []string{
		"../tests/p1.peepoo",
		"../tests/p2.peepoo",
		"../tests/p3.peepoo",
		"../tests/p4.peepoo",
	}

	tokenList := "../tests/tokens.list"

	for _, progPath := range progList {
		_, tokens := Tokenize(tokenList, progPath)
		if !compareTokens(tokens, &expected) {
			t.Errorf("Tokens not equal for tokens %s and program %s.\nActual:\n%s\nExpected:\n%s", tokenList, progPath, *tokens, expected)
		}
	}

}

func TestTokenizerSeperator(t *testing.T) {
	fmt.Println(os.Getwd())
	expected := []Token{
		{Name: "integer"},
		{Name: "plus"},
		{Name: "integer"},
		{Name: "seperator"},
		{Name: "integer"},
		{Name: "plus"},
		{Name: "integer"},
		{Name: "seperator"},
	}

	progList := []string{
		"../tests/p5.peepoo",
	}

	tokenList := "../tests/tokens.list"

	for _, progPath := range progList {
		_, tokens := Tokenize(tokenList, progPath)
		if !compareTokens(tokens, &expected) {
			t.Errorf("Tokens not equal for tokens %s and program %s.\nActual:\n%s\nExpected:\n%s", tokenList, progPath, *tokens, expected)
		}
	}

}
