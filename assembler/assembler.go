package assembler

import (
	"JureBevc/gpc/parser"
	"JureBevc/gpc/util"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type AssemblerRule struct {
	ParentName    string
	ChildrenNames []string
	Template      string
}

func loadAssemblerFile(assemblerFile string) *[]AssemblerRule {
	file, err := os.Open(assemblerFile)
	if err != nil {
		log.Fatalf("Unable to open assembler file with path %s\n%s\n", assemblerFile, err)
		return nil
	}

	defer file.Close()

	rules := []AssemblerRule{}
	scanner := bufio.NewScanner(file)
	currentRule := AssemblerRule{}
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "---" {
			if currentRule.ParentName != "" {
				rules = append(rules, currentRule)
			}
			currentRule = AssemblerRule{}
			continue
		}

		if currentRule.ParentName == "" {
			currentRule.ParentName = line
		} else if len(currentRule.ChildrenNames) == 0 {
			lineSplit := strings.Split(line, " ")
			for _, symbol := range lineSplit {
				symbolTrimmed := strings.TrimSpace(symbol)
				if symbolTrimmed != "" {
					currentRule.ChildrenNames = append(currentRule.ChildrenNames, symbolTrimmed)
				}
			}
		} else {
			if currentRule.Template == "" {
				currentRule.Template = line
			} else {
				currentRule.Template = currentRule.Template + "\n" + line
			}
		}
	}

	return &rules
}

func assembleNode(node *util.TreeNode[parser.ParseNode], rules *[]AssemblerRule, outputFile *os.File) string {
	// Terminals are leafs, return their value
	if node.Value.IsTerminal {
		return node.Value.Value
	}

	// Check if any rules apply
	for _, rule := range *rules {
		if rule.ParentName == node.Value.Name && len(rule.ChildrenNames) == len(node.Children) {
			ruleMatch := true
			for childIndex, childName := range rule.ChildrenNames {
				if childName != node.Children[childIndex].Value.Name {
					ruleMatch = false
					break
				}
			}

			if ruleMatch {
				template := rule.Template

				for _, childNode := range node.Children {
					childAssembly := assembleNode(childNode, rules, outputFile)
					tagName := "$" + childNode.Value.Name + "$"
					template = strings.ReplaceAll(template, tagName, childAssembly)
				}

				return template
			}
		}
	}

	// No rules apply, just process children
	combinedAssembly := ""
	for _, childNode := range node.Children {
		childAssembly := assembleNode(childNode, rules, outputFile)
		combinedAssembly = combinedAssembly + childAssembly
	}

	return combinedAssembly
}

func assembleTree(parseTree *util.TreeNode[parser.ParseNode], rules *[]AssemblerRule) {
	file, err := os.Create("out.gpc")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function exits

	assembly := assembleNode(parseTree, rules, file)

	file.Write([]byte(assembly))
}

func Assemble(parseTree *util.TreeNode[parser.ParseNode], assemblerFile string) {
	rules := loadAssemblerFile(assemblerFile)
	assembleTree(parseTree, rules)
}
