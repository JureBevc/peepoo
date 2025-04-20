package runtime

import (
	"JureBevc/peepoo/parser"
	"JureBevc/peepoo/util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Scope map[string]interface{}

var temporaryVars int = 0

func CopyScope(scope *Scope) *Scope {
	newScope := Scope{}
	for key, val := range *scope {
		newScope[key] = val
	}
	return &newScope
}

func ScopeIsReturning(scope *Scope) bool {
	if _, ok := (*scope)["RET"]; ok {
		return true
	}

	return false
}

func RunAssign(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	variableName := node.Children[0].Value.Value
	valueNode := node.Children[2]
	result := RunMath(valueNode, scope)
	(*scope)[variableName] = result
}

func RunValue(node *util.TreeNode[parser.ParseNode], scope *Scope) interface{} {
	if node.Value.Name == "VALUE" {
		firstChild := node.Children[0]
		switch firstChild.Value.Name {
		case "var":
			if varValue, ok := (*scope)[firstChild.Value.Value]; ok {
				return varValue
			} else {
				log.Fatalf("Undefined variable %s\n", firstChild.Value.Value)
			}
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
	if len(node.Children) == 1 && node.Children[0].Value.Name == "VALUE" {
		valueChild := node.Children[0]
		val := RunValue(valueChild, scope)
		return val
	}

	if len(node.Children) == 1 && node.Children[0].Value.Name == "FUNCCALL" {
		valueChild := node.Children[0]
		val := RunFuncCall(valueChild, scope)
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

				newRightNode := parser.CopyTree(rightMathNode)

				newRightNode.Children[0].Children[0].Value.Name = "var"
				newRightNode.Children[0].Children[0].Value.Value = tmpVarName
				(*scope)[tmpVarName] = newValue

				return RunMath(newRightNode, scope)
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
			if ScopeIsReturning(scope) {
				break
			}
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
			if ScopeIsReturning(scope) {
				break
			}
			bodyNode = bodyNode.Children[1]
		}
		currentValue += 1
	}
}

func RunFunc(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	funcVariableName := node.Children[1].Value.Value
	funcParamNode := node.Children[2]
	(*scope)[funcVariableName] = funcParamNode
}

func RunReturn(node *util.TreeNode[parser.ParseNode], scope *Scope) {
	if len(node.Children) == 1 {
		(*scope)["RET"] = nil
	}

	mathNode := node.Children[1]
	value := RunMath(mathNode, scope)
	(*scope)["RET"] = value
}

func RunFuncCall(node *util.TreeNode[parser.ParseNode], scope *Scope) interface{} {
	funcVariableName := node.Children[1].Value.Value
	funcParamNode := (*scope)[funcVariableName].(*util.TreeNode[parser.ParseNode])

	funcParamNames := []string{}
	for funcParamNode.Children[1].Value.Name != "FUNCBODY" {
		varName := funcParamNode.Children[0].Value.Value
		funcParamNames = append(funcParamNames, varName)
		funcParamNode = funcParamNode.Children[1]
	}

	funcBodyNode := funcParamNode.Children[1]

	callParamNode := node.Children[2]
	mathNodes := []*util.TreeNode[parser.ParseNode]{}
	for callParamNode.Children[0].Value.Name != "funccall" {
		mathNodes = append(mathNodes, callParamNode.Children[0])
		callParamNode = callParamNode.Children[1]
	}

	if len(mathNodes) != len(funcParamNames) {
		log.Fatalf("Unmatching number of function paramaters. %s expected %d, got %d.\n", funcVariableName, len(funcParamNames), len(mathNodes))
	}

	scopeCopy := CopyScope(scope)

	for i := 0; i < len(mathNodes); i++ {
		value := RunMath(mathNodes[i], scope)
		varName := funcParamNames[i]
		(*scopeCopy)[varName] = value
	}

	for funcBodyNode.Children[0].Value.Name != "funcend" {
		expressionNode := funcBodyNode.Children[0]
		RunExpression(expressionNode, scopeCopy)
		if returnValue, ok := (*scopeCopy)["RET"]; ok {
			return returnValue
		}
		funcBodyNode = funcBodyNode.Children[1]
	}

	return nil
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
		case "FUNC":
			RunFunc(childNode, scope)
		case "FUNCRETURN":
			RunReturn(childNode, scope)
		case "FUNCCALL":
			RunFuncCall(childNode, scope)
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
			if ScopeIsReturning(scope) {
				break
			}
		}

		currentProgram = nextProgram
	}
}

func RunTree(parseTree *util.TreeNode[parser.ParseNode]) {
	newScope := Scope{}
	RunProgram(parseTree, &newScope)
}
