package parser

import (
	"JureBevc/gpc/tokenizer"
	"JureBevc/gpc/util"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type GrammarSymbol struct {
	Name       string
	IsTerminal bool
}

// Map non-terminal name to list of rules (where every rule is a list of symbols)
type GrammarRules map[string][][]GrammarSymbol

func stringIsTerminal(name string, allTerminals *[]tokenizer.TokenDefinition) bool {
	isTerminal := false
	for _, definition := range *allTerminals {
		if definition.Name == name {
			isTerminal = true
			break
		}
	}

	return isTerminal
}

func loadGrammarFile(pathToGrammarFile string, allTerminals *[]tokenizer.TokenDefinition) (*GrammarRules, GrammarSymbol) {
	file, err := os.Open(pathToGrammarFile)
	if err != nil {
		log.Fatalf("Unable to open grammar file with path %s\n%s\n", pathToGrammarFile, err)
		return nil, GrammarSymbol{}
	}

	defer file.Close()

	grammar := GrammarRules{}
	firstSymbol := GrammarSymbol{}
	scanner := bufio.NewScanner(file)
	currentNonTerminal := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			currentNonTerminal = ""
			continue
		}

		if currentNonTerminal == "" {
			// New non-terminal entry
			currentNonTerminal = line
			var newEntry [][]GrammarSymbol
			grammar[currentNonTerminal] = newEntry
			if firstSymbol.Name == "" {
				firstSymbol = GrammarSymbol{Name: currentNonTerminal, IsTerminal: false}
			}
		} else {
			// New rule for current non-terminal
			rule := strings.Split(line, " ")
			var newRule []GrammarSymbol
			for _, name := range rule {
				grammarSymbol := GrammarSymbol{
					Name:       name,
					IsTerminal: stringIsTerminal(name, allTerminals),
				}
				newRule = append(newRule, grammarSymbol)
			}
			grammar[currentNonTerminal] = append(grammar[currentNonTerminal], newRule)
		}
	}

	// Validation
	for key := range grammar {
		rules := grammar[key]
		for _, rule := range rules {
			for _, symbol := range rule {
				// Every symbol must be a terminal or non-terminal

				// Check for non-terminal
				_, isNonTerminal := grammar[symbol.Name]

				if isNonTerminal {
					continue
				}

				// Check for terminal
				isTerminal := false
				for _, definition := range *allTerminals {
					if definition.Name == symbol.Name {
						isTerminal = true
						break
					}
				}

				if isTerminal {
					continue
				} else {
					log.Panicf("Unknown symbol in grammar: %s\n", symbol.Name)
				}
			}
		}
	}

	return &grammar, firstSymbol
}

func naiveParseRecursive(programTokens *[]tokenizer.Token, grammar *GrammarRules, currentSymbol GrammarSymbol, startSymbol GrammarSymbol, tokenIndex int) (*util.TreeNode[string], int) {
	// Terminals have no rules, return as leaf node
	if currentSymbol.IsTerminal {

		if currentSymbol.Name != (*programTokens)[tokenIndex].Name {
			// Terminal cannot match
			return nil, tokenIndex
		}

		// Terminal can match
		return &util.TreeNode[string]{
			Children: nil,
			Value:    currentSymbol.Name,
		}, tokenIndex + 1
	}

	// Loop non-terminal rules and try to parse each one
	rules := (*grammar)[currentSymbol.Name]

	for _, rule := range rules {
		var children []*util.TreeNode[string]

		parsedAllChildren := true
		childTokenIndex := tokenIndex
		for _, childSymbol := range rule {
			var childNode *util.TreeNode[string]
			childNode, childTokenIndex = naiveParseRecursive(programTokens, grammar, childSymbol, startSymbol, childTokenIndex)
			if childNode == nil {
				// Could not create children, rule cannot apply
				parsedAllChildren = false
				break
			} else {
				children = append(children, childNode)
			}
		}

		if parsedAllChildren && currentSymbol.Name == startSymbol.Name {
			// Start symbol must also match end of file
			if childTokenIndex != len(*programTokens) {
				parsedAllChildren = false
			}
		}

		// Parsing children was a success, return result
		if parsedAllChildren {
			return &util.TreeNode[string]{
				Children: children,
				Value:    currentSymbol.Name,
			}, childTokenIndex
		}

	}

	return nil, tokenIndex
}

func naiveParse(programTokens *[]tokenizer.Token, grammar *GrammarRules, firstSymbol GrammarSymbol) *util.TreeNode[string] {
	tree, _ := naiveParseRecursive(programTokens, grammar, firstSymbol, firstSymbol, 0)
	return tree
}

func Parse(terminals *[]tokenizer.TokenDefinition, programTokens *[]tokenizer.Token, grammarFile string) *GrammarRules {
	grammar, firstSymbol := loadGrammarFile(grammarFile, terminals)
	parseTree := naiveParse(programTokens, grammar, firstSymbol)
	if parseTree == nil {
		fmt.Println(fmt.Errorf("could not create parse tree"))
	} else {
		util.PrintTree(parseTree, "-")
	}
	return grammar
}
