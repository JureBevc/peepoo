package runtime

import (
	"JureBevc/poopoo/parser"
	"JureBevc/poopoo/util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Scope map[string]interface{}

var temporaryVars int = 0

func RunAssign(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	variableName := node.Children[0].Value.Value
	valueNode := node.Children[2]
	result := RunMath(valueNode, scope)
	(*scope)[variableName] = result
}

func RunOpMath(node *util.TreeNode[parser.ParseNode], scope *Scope) interface{} {
	secondChild := node.Children[1]
	switch secondChild.Value.Name {
	case "plus":
		result1 := RunMath(node.Children[0], scope).(int64)
		result2 := RunMath(node.Children[2], scope).(int64)
		return result1 + result2
	case "minus":
		result1 := RunMath(node.Children[0], scope).(int64)
		result2 := RunMath(node.Children[2], scope).(int64)
		return result1 - result2
	}

	log.Fatalf("Failed to parse operator %v\n", node)
	return nil
}

func RunValue(node *util.TreeNode[parser.ParseNode], scope *Scope) interface{} {
	if node.Value.Name == "VALUE" {
		firstChild := node.Children[0]
		switch firstChild.Value.Name {
		case "var":
			return (*scope)[firstChild.Value.Value]
		case "binary":
			binaryStr := strings.ReplaceAll(firstChild.Value.Value, "p", "")
			binaryStr = strings.ReplaceAll(strings.ReplaceAll(binaryStr, "i", "1"), "o", "0")
			val, err := strconv.ParseInt(binaryStr, 2, 64)
			if err != nil {
				log.Fatalf("Failed to parse binary number from %s\n", firstChild.Value.Value)
			}
			return val
		}
	}

	log.Fatalf("Failed to parse value %v\n", node)
	return nil
}

func RunMath(node *util.TreeNode[parser.ParseNode], scope *Scope) interface{} {
	if len(node.Children) == 1 {
		valueChild := node.Children[0]
		val := RunValue(valueChild, scope)
		return val
	}

	if len(node.Children) == 2 {
		if node.Children[0].Value.Name == "VALUE" &&
			node.Children[1].Value.Name == "OP_MATH" {
			leftValue := RunValue(node.Children[0], scope).(int64)

			operator := node.Children[1].Children[0].Value.Name
			rightMathNode := node.Children[1].Children[1]
			switch operator {
			case "plus":
				rightValue := RunMath(rightMathNode, scope).(int64)
				return leftValue + rightValue
			case "minus":
				rightValue := RunMath(rightMathNode, scope).(int64)
				return leftValue - rightValue
			case "multiply":
				rightValue := RunValue(rightMathNode.Children[0], scope).(int64)
				newValue := leftValue * rightValue

				temporaryVars += 1
				tmpVarName := fmt.Sprintf("TMP%d", temporaryVars)
				rightMathNode.Children[0].Children[0].Value.Name = "var"
				rightMathNode.Children[0].Children[0].Value.Value = tmpVarName
				(*scope)[tmpVarName] = newValue

				return RunMath(rightMathNode, scope)
			case "divide":
				rightValue := RunValue(rightMathNode.Children[0], scope).(int64)
				newValue := leftValue / rightValue

				temporaryVars += 1
				tmpVarName := fmt.Sprintf("TMP%d", temporaryVars)
				rightMathNode.Children[0].Children[0].Value.Name = "var"
				rightMathNode.Children[0].Children[0].Value.Value = tmpVarName
				(*scope)[tmpVarName] = newValue

				return RunMath(rightMathNode, scope)
			}
		}
	}

	log.Fatalf("Failed to run math expression %v %d\n", node.Value, len(node.Children))
	return nil
}

func RunIf(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	mathNode := node.Children[1]
	mathValue := RunMath(mathNode, scope).(int64)
	if mathValue != 0 {
		bodyNode := node.Children[2]

		// body node has 1 child when its ifend
		for len(bodyNode.Children) > 1 {
			RunExpression(bodyNode.Children[0], scope)
			bodyNode = bodyNode.Children[1]
		}

	}
}

func RunLoop(node *util.TreeNode[parser.ParseNode], scope *Scope) {

	varNode := node.Children[1]
	variableName := varNode.Value.Value

	mathValueStart := RunMath(node.Children[2], scope).(int64)
	mathValueStop := RunMath(node.Children[3], scope).(int64)
	currentValue := mathValueStart
	for currentValue < mathValueStop {
		(*scope)[variableName] = currentValue
		bodyNode := node.Children[4]

		// body node has 1 child when its ifend
		for len(bodyNode.Children) > 1 {
			RunExpression(bodyNode.Children[0], scope)
			bodyNode = bodyNode.Children[1]
		}
		currentValue += 1
	}
}

func RunPrint(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	result := RunMath(node.Children[1], scope)
	fmt.Print(result)
}

func RunPrintln(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	result := RunMath(node.Children[1], scope)
	fmt.Println(result)
}

func RunExpression(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	for _, childNode := range node.Children {
		switch childNode.Value.Name {
		case "ASSIGN":
			RunAssign(childNode, scope)
		case "PRINT":
			RunPrint(childNode, scope)
		case "PRINTLN":
			RunPrintln(childNode, scope)
		case "IF":
			RunIf(childNode, scope)
		case "LOOP":
			RunLoop(childNode, scope)
		}
	}
}

func RunProgram(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	if node.Value.Name != "PROGRAM" {
		log.Fatalf("Failed to run program, unexpected node %s\n", node.Value.Name)
	}

	currentProgram := node
	for currentProgram != nil {
		var nextProgram *util.TreeNode[parser.ParseNode] = nil
		var expressionNode *util.TreeNode[parser.ParseNode] = nil

		for _, childNode := range currentProgram.Children {
			if childNode.Value.Name == "EXPRESSION" {
				expressionNode = childNode
			} else if childNode.Value.Name == "PROGRAM" {
				nextProgram = childNode
			}
		}

		if expressionNode != nil {
			RunExpression(expressionNode, scope)
		}

		currentProgram = nextProgram
	}
}

func RunTree(parseTree *util.TreeNode[parser.ParseNode]) {
	newScope := Scope{}
	RunProgram(parseTree, &newScope)
}
