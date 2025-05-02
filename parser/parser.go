package parser

import (
	"JureBevc/peepoo/tokenizer"
	"JureBevc/peepoo/util"
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"log"
	"strings"
)

type GrammarSymbol struct {
	Name       string
	IsTerminal bool
}

type ParseNode struct {
	Name       string
	Value      string
	IsTerminal bool
	Token      *tokenizer.Token
}

var MaxLine int = 0
var MaxColumn int = 0

func CopyTree(node *util.TreeNode[ParseNode]) *util.TreeNode[ParseNode] {
	if len((*node).Children) == 0 {
		newNode := util.TreeNode[ParseNode]{
			Children: nil,
			Value:    node.Value,
		}
		return &newNode
	}

	newChildren := []*util.TreeNode[ParseNode]{}

	for _, cn := range (*node).Children {
		newChildren = append(newChildren, CopyTree(cn))
	}

	newNode := util.TreeNode[ParseNode]{
		Children: newChildren,
		Value:    node.Value,
	}
	return &newNode
}

// Map non-terminal name to list of rules (where every rule is a list of symbols)
type GrammarRules map[string][][]GrammarSymbol

func PrintTree(tree *util.TreeNode[ParseNode], prefix string) {
	fmt.Printf("%s%s (%s)\n", prefix, fmt.Sprint(tree.Value.Value), tree.Value.Name)
	childPrefix := prefix + "|"
	for _, child := range (*tree).Children {
		PrintTree(child, childPrefix)
	}
}

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

func loadGrammarFile(pathToGrammarFile embed.FS, allTerminals *[]tokenizer.TokenDefinition) (*GrammarRules, GrammarSymbol) {
	file, err := pathToGrammarFile.ReadFile("config/grammar.list")
	if err != nil {
		log.Fatalf("Unable to open grammar file with path %v\n%s\n", pathToGrammarFile, err)
		return nil, GrammarSymbol{}
	}

	grammar := GrammarRules{}
	firstSymbol := GrammarSymbol{}
	scanner := bufio.NewScanner(bytes.NewReader(file))
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

func updateMaxLineAndColumn(line int, column int) {
	if line > MaxLine {
		MaxLine = line
		MaxColumn = column
	} else if line == MaxLine && column > MaxColumn {
		MaxColumn = column
	}
}

func naiveParseRecursive(programTokens *[]tokenizer.Token, grammar *GrammarRules, currentSymbol GrammarSymbol, startSymbol GrammarSymbol, tokenIndex int) (*util.TreeNode[ParseNode], int) {
	// Terminals have no rules, return as leaf node
	if tokenIndex >= len(*programTokens) {
		return nil, tokenIndex
	}

	currentToken := (*programTokens)[tokenIndex]
	if currentSymbol.IsTerminal {
		if currentSymbol.Name != currentToken.Name {
			// Terminal cannot match
			return nil, tokenIndex
		}

		updateMaxLineAndColumn(currentToken.Line, currentToken.Column)
		// Terminal can match
		return &util.TreeNode[ParseNode]{
			Children: nil,
			Value: ParseNode{
				Name:       currentToken.Name,
				Value:      currentToken.Value,
				IsTerminal: true,
				Token:      &currentToken,
			},
		}, tokenIndex + 1
	}

	// Loop non-terminal rules and try to parse each one
	rules := (*grammar)[currentSymbol.Name]
	for _, rule := range rules {
		var children []*util.TreeNode[ParseNode]
		parsedAllChildren := true
		childTokenIndex := tokenIndex
		for _, childSymbol := range rule {
			var childNode *util.TreeNode[ParseNode]
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
			return &util.TreeNode[ParseNode]{
				Children: children,
				Value: ParseNode{
					Name:       currentSymbol.Name,
					Value:      currentSymbol.Name,
					IsTerminal: false,
					Token:      &currentToken,
				},
			}, childTokenIndex
		}

	}

	return nil, tokenIndex
}

func naiveParse(programTokens *[]tokenizer.Token, grammar *GrammarRules, firstSymbol GrammarSymbol) *util.TreeNode[ParseNode] {
	tree, _ := naiveParseRecursive(programTokens, grammar, firstSymbol, firstSymbol, 0)
	if tree == nil {
		util.FatalError("Failed to parse expression", MaxLine, MaxColumn)
	}
	return tree
}

func Parse(terminals *[]tokenizer.TokenDefinition, programTokens *[]tokenizer.Token, grammarFile embed.FS) *util.TreeNode[ParseNode] {
	grammar, firstSymbol := loadGrammarFile(grammarFile, terminals)
	parseTree := naiveParse(programTokens, grammar, firstSymbol)
	return parseTree
}
