package tokenizer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type TokenDefinition struct {
	Name       string
	Definition string
	Regex      *regexp.Regexp
	IsRegex    bool
}

type Token struct {
	Name  string
	Value string
}

func loadTokenFile(pathToTokenFile string) *[]TokenDefinition {
	file, err := os.Open(pathToTokenFile)
	if err != nil {
		log.Fatalf("Unable to open token file with path %s\n%s\n", pathToTokenFile, err)
		return nil
	}

	defer file.Close()

	var tokens []TokenDefinition

	scanner := bufio.NewScanner(file)
	currentDefinition := TokenDefinition{}
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if currentDefinition.Name == "" {
			currentDefinition.Name = line
		} else {
			definitionString, hasPrefix := strings.CutPrefix(line, "regex:")
			currentDefinition.Definition = definitionString
			currentDefinition.IsRegex = hasPrefix

			if hasPrefix {
				currentDefinition.Regex = regexp.MustCompile(currentDefinition.Definition)
			}

			tokens = append(tokens, currentDefinition)
			currentDefinition = TokenDefinition{}
		}
	}

	return &tokens
}

func wordSingleDefinition(tokenDefinitions *[]TokenDefinition, word string) (TokenDefinition, error) {
	// Returns a single token definition, if there is only one definition that is valid
	// If zero or more than one definitons exist, it return an error

	validDefinition := TokenDefinition{}
	validDefinitionFound := false
	for _, definition := range *tokenDefinitions {
		isValid := false
		if definition.IsRegex {
			isValid = definition.Regex.MatchString(word)
		} else {
			isValid = definition.Definition == word
		}

		if isValid && validDefinitionFound {
			return TokenDefinition{}, fmt.Errorf("more than one definition found for word %s", word)
		}

		if isValid {
			validDefinition = definition
			validDefinitionFound = true
		}
	}

	if !validDefinitionFound {
		return TokenDefinition{}, fmt.Errorf("no definition found for word %s", word)
	}

	return validDefinition, nil
}

func parseFile(tokenDefinitons *[]TokenDefinition, pathToInputFile string) *[]Token {
	var tokens []Token

	file, err := os.Open(pathToInputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	currentWord := ""
	for {
		char, _, err := reader.ReadRune()
		reachedEnd := false
		// Check for end of file
		if err != nil {
			reachedEnd = true
		}

		// Create next word if current char is not empty
		emptyChar := unicode.IsSpace(char)
		nextWord := currentWord
		if !reachedEnd && !emptyChar {
			nextWord = currentWord + string(char)
		}

		// Check if current and next word are parsable
		currentDefinition, currentDefError := wordSingleDefinition(tokenDefinitons, currentWord)
		_, nextDefError := wordSingleDefinition(tokenDefinitons, nextWord)

		parseCurrentWord := false
		if emptyChar && currentWord != "" {
			parseCurrentWord = true
		}

		if currentDefError == nil && nextDefError != nil {
			// Current word is a token, next word is not, create a split
			parseCurrentWord = true
		}

		if reachedEnd && currentWord != "" {
			// Last word should parse
			parseCurrentWord = true
		}

		if parseCurrentWord {
			if currentDefError != nil {
				log.Fatalf("Unable to parse token word %s\n", currentWord)
			}
			tokens = append(tokens, Token{
				Name:  currentDefinition.Name,
				Value: currentWord,
			})
		}

		if reachedEnd {
			break
		}

		if parseCurrentWord {
			currentWord = ""
			if !emptyChar {
				currentWord = string(char)
			}
		} else {
			currentWord = nextWord
		}
	}

	return &tokens
}

func Tokenize(pathToTokenFile string, pathToInputFile string) (*[]TokenDefinition, *[]Token) {
	tokenDef := loadTokenFile(pathToTokenFile)
	tokens := parseFile(tokenDef, pathToInputFile)
	return tokenDef, tokens
}
