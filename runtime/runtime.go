package runtime

import (
	"JureBevc/peepoo/parser"
	"JureBevc/peepoo/util"
	"bufio"
	"fmt"
	"log"
	"os"
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

func EncodeString(s string) string {
	var base5Map = []string{"pa", "pe", "pi", "po", "pu"}
	var builder strings.Builder
	for i, c := range []byte(s) {
		val := int(c)
		var digits [4]int

		// Convert to base 5, right-aligned to 4 digits
		for j := 3; j >= 0; j-- {
			digits[j] = val % 5
			val /= 5
		}

		for _, d := range digits {
			builder.WriteString(base5Map[d])
		}

		if i < len(s)-1 {
			builder.WriteByte(' ')
		}
	}

	return builder.String()
}

func DecodeString(encoded string) (string, error) {
	revMap := map[string]int{
		"pa": 0, "pe": 1, "pi": 2, "po": 3, "pu": 4,
	}

	words := strings.Split(encoded, " ")
	var result strings.Builder

	for _, word := range words {
		if len(word) != 8 {
			return "", fmt.Errorf("failed to decode word %s", word)
		}

		val := 0
		for i := 0; i < 8; i += 2 {
			symbol := word[i : i+2]
			digit, ok := revMap[symbol]
			if !ok {
				return "", fmt.Errorf("failed to decode symbol %s", symbol)
			}
			val = val*5 + digit
		}

		err := result.WriteByte(byte(val))
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func RunAssign(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	switch node.Children[0].Value.Name {
	case "var":
		variableName := node.Children[0].Value.Value
		valueNode := node.Children[2]
		result, err := RunMath(valueNode, scope)
		if err != nil {
			return err
		}
		(*scope)[variableName] = result
	case "LISTACCESS":
		variableName := node.Children[0].Value.Value
		valueNode := node.Children[2]
		result, err := RunMath(valueNode, scope)
		if err != nil {
			return err
		}
		accessNode := node.Children[0]
		if varValue, ok := (*scope)[accessNode.Children[0].Value.Value]; ok {
			varList := varValue.([]interface{})
			indexValue, err := RunValue(accessNode.Children[2], scope)
			if err != nil {
				return err
			}
			indexInt, ok := indexValue.(int64)
			if !ok {
				return util.FormatError(
					"List index not an integer",
					accessNode.Value.Token.Line,
					accessNode.Value.Token.Column,
				)
			}
			varList[indexInt] = result
			(*scope)[variableName] = varList
		}
	}
	return nil
}

func RunValue(node *util.TreeNode[parser.ParseNode], scope *Scope) (interface{}, error) {
	if node.Value.Name == "VALUE" {
		firstChild := node.Children[0]
		switch firstChild.Value.Name {
		case "var":
			if varValue, ok := (*scope)[firstChild.Value.Value]; ok {
				return varValue, nil
			} else {
				util.FatalError(
					fmt.Sprintf("Undefined variable %s", firstChild.Value.Value),
					firstChild.Value.Token.Line,
					firstChild.Value.Token.Column,
				)
			}
		case "binary":
			binaryStr := strings.ReplaceAll(firstChild.Value.Value, "p", "")
			binaryStr = strings.ReplaceAll(strings.ReplaceAll(binaryStr, "i", "1"), "o", "0")
			val, err := strconv.ParseInt(binaryStr, 2, 64)
			if err != nil {
				util.FatalError(
					fmt.Sprintf("Failed to parse binary number from %s", firstChild.Value.Value),
					firstChild.Value.Token.Line,
					firstChild.Value.Token.Column,
				)
			}
			return val, nil
		case "char":
			return DecodeString(firstChild.Value.Value)
		case "FUNCCALL":
			return RunFuncCall(firstChild, scope)
		case "LIST":
			return ParseList(firstChild, scope)
		case "LISTACCESS":
			if varValue, ok := (*scope)[firstChild.Children[0].Value.Value]; ok {
				varList := varValue.([]interface{})
				indexValue, err := RunValue(firstChild.Children[2], scope)
				if err != nil {
					return nil, err
				}
				indexInt, ok := indexValue.(int64)
				if !ok {
					return nil, util.FormatError(
						"List index is not an integer",
						firstChild.Value.Token.Line,
						firstChild.Value.Token.Column,
					)
				}
				listLen := int64(len(varList))
				if indexInt < 0 || indexInt >= listLen {
					errorText := fmt.Sprintf("List index %d out of range [%d, %d]", indexInt, 0, listLen-1)
					return nil, util.FormatError(errorText, node.Value.Token.Line, node.Value.Token.Column)
				}
				return varList[indexInt], nil
			}
			errorText := fmt.Sprintf("Failed to access list %s", firstChild.Value.Value)
			util.FatalError(errorText, node.Value.Token.Line, node.Value.Token.Column)
		case "LISTPOP":
			return RunListPop(firstChild, scope)
		case "LISTLEN":
			if varValue, ok := (*scope)[firstChild.Children[1].Value.Value]; ok {
				varList, ok := varValue.([]interface{})
				if !ok {
					return nil, util.FormatError(
						"Failed to parse list to get list length",
						firstChild.Value.Token.Line,
						firstChild.Value.Token.Column,
					)
				}
				return int64(len(varList)), nil
			}
		case "readinput":
			reader := bufio.NewReader(os.Stdin)
			data, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			data = strings.TrimSpace(data)
			var chars []interface{}
			for _, r := range data {
				chars = append(chars, string(r))
			}
			return chars, nil
		case "readfile":
			if varValue, ok := (*scope)[node.Children[1].Value.Value]; ok {
				listValue := varValue.([]interface{})
				varString := ""
				for _, val := range listValue {
					varString = varString + val.(string)
				}
				data, _ := os.ReadFile(varString)
				outString := string(data)
				var chars []interface{}
				for _, r := range outString {
					chars = append(chars, string(r))
				}

				return chars, nil
			}
			errorText := fmt.Sprintf("Failed to read file %s\n", firstChild.Children[1].Value.Value)
			util.FatalError(errorText, node.Value.Token.Line, node.Value.Token.Column)
		case "chartoint":
			secondChild := node.Children[1]
			if secondChild.Value.Name == "var" {
				if varValue, ok := (*scope)[secondChild.Value.Value]; ok {
					charVal := varValue.(string)
					return int64(rune(charVal[0])), nil
				}
			}
			if secondChild.Value.Name == "char" {
				charVal, err := DecodeString(secondChild.Value.Value)
				if err != nil {
					return nil, err
				}
				return int64(rune(charVal[0])), nil
			}
		}

		return nil, util.FormatError(
			"Failed to parse value",
			firstChild.Value.Token.Line,
			firstChild.Value.Token.Column,
		)
	}

	return nil, util.FormatError(
		"Failed to parse value",
		node.Value.Token.Line,
		node.Value.Token.Column,
	)
}

func ParseList(node *util.TreeNode[parser.ParseNode], scope *Scope) ([]interface{}, error) {
	ret := []interface{}{}

	listElement := node.Children[1]
	for listElement.Value.Name == "LISTELEMENT" {
		if len(listElement.Children) > 1 {
			val, err := RunValue(listElement.Children[0], scope)
			if err != nil {
				return nil, err
			}
			ret = append(ret, val)
			listElement = listElement.Children[1]
		} else {
			listElement = listElement.Children[0]
		}
	}

	return ret, nil
}

func RunMath(node *util.TreeNode[parser.ParseNode], scope *Scope) (interface{}, error) {
	if len(node.Children) == 1 && node.Children[0].Value.Name == "VALUE" {
		valueChild := node.Children[0]
		return RunValue(valueChild, scope)
	}

	if len(node.Children) == 1 && node.Children[0].Value.Name == "FUNCCALL" {
		valueChild := node.Children[0]
		return RunFuncCall(valueChild, scope)
	}

	if len(node.Children) == 2 {
		if node.Children[0].Value.Name == "VALUE" &&
			node.Children[1].Value.Name == "OP_MATH" {
			result, err := RunValue(node.Children[0], scope)
			if err != nil {
				return nil, err
			}
			leftValue, ok := result.(int64)
			if !ok {
				util.FatalError("Invalid value type", node.Children[0].Value.Token.Line, node.Children[0].Value.Token.Column)
			}

			operator := node.Children[1].Children[0].Value.Name
			rightMathNode := node.Children[1].Children[1]
			switch operator {
			case "plus":
				val, err := RunMath(rightMathNode, scope)
				if err != nil {
					return nil, err
				}
				rightValue, ok := val.(int64)
				if !ok {
					util.FatalError(
						"Invalid value type",
						node.Children[0].Value.Token.Line,
						node.Children[0].Value.Token.Column,
					)
				}
				return leftValue + rightValue, nil
			case "minus":
				val, err := RunMath(rightMathNode, scope)
				if err != nil {
					return nil, err
				}
				rightValue, ok := val.(int64)
				if !ok {
					util.FatalError(
						"Invalid value type",
						node.Children[0].Value.Token.Line,
						node.Children[0].Value.Token.Column,
					)
				}
				return leftValue - rightValue, nil
			case "multiply":
				val, err := RunValue(rightMathNode.Children[0], scope)
				if err != nil {
					return nil, err
				}
				rightValue, ok := val.(int64)
				if !ok {
					util.FatalError(
						"Invalid value type",
						node.Children[0].Value.Token.Line,
						node.Children[0].Value.Token.Column,
					)
				}
				newValue := leftValue * rightValue

				temporaryVars += 1
				tmpVarName := fmt.Sprintf("TMP%d", temporaryVars)

				newRightNode := parser.CopyTree(rightMathNode)

				newRightNode.Children[0].Children[0].Value.Name = "var"
				newRightNode.Children[0].Children[0].Value.Value = tmpVarName
				(*scope)[tmpVarName] = newValue

				return RunMath(newRightNode, scope)
			case "divide":
				val, err := RunValue(rightMathNode.Children[0], scope)
				if err != nil {
					return nil, err
				}
				rightValue, ok := val.(int64)
				if !ok {
					util.FatalError(
						"Invalid value type",
						node.Children[0].Value.Token.Line,
						node.Children[0].Value.Token.Column,
					)
				}
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

	return nil, util.FormatError(
		"Failed to run math expression",
		node.Value.Token.Line,
		node.Value.Token.Column,
	)
}

func RunIf(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	mathNode := node.Children[1]
	val, err := RunMath(mathNode, scope)
	if err != nil {
		return err
	}
	mathValue, ok := val.(int64)
	if !ok {
		return util.FormatError(
			"Invalid value type",
			node.Value.Token.Line,
			node.Value.Token.Column,
		)
	}
	if mathValue != 0 {
		bodyNode := node.Children[2]

		// body node has 1 child when its ifend
		for len(bodyNode.Children) > 1 {
			err := RunExpression(bodyNode.Children[0], scope)
			if err != nil {
				return err
			}
			if ScopeIsReturning(scope) {
				break
			}
			bodyNode = bodyNode.Children[1]
		}

	}
	return nil
}

func RunLoop(node *util.TreeNode[parser.ParseNode], scope *Scope) error {

	varNode := node.Children[1]
	variableName := varNode.Value.Value

	startValue, err := RunMath(node.Children[2], scope)
	if err != nil {
		return err
	}

	stopValue, err := RunMath(node.Children[3], scope)
	if err != nil {
		return err
	}

	mathValueStart, ok := startValue.(int64)
	if !ok {
		util.FatalError(
			"Invalid value type for loop start value",
			node.Children[0].Value.Token.Line,
			node.Children[0].Value.Token.Column,
		)
	}
	mathValueStop, ok := stopValue.(int64)
	if !ok {
		util.FatalError(
			"Invalid value type for loop stop value",
			node.Children[0].Value.Token.Line,
			node.Children[0].Value.Token.Column,
		)
	}

	currentValue := mathValueStart
	for currentValue < mathValueStop {
		(*scope)[variableName] = currentValue
		bodyNode := node.Children[4]

		// body node has 1 child when its ifend
		for len(bodyNode.Children) > 1 {
			err := RunExpression(bodyNode.Children[0], scope)
			if err != nil {
				return err
			}
			if ScopeIsReturning(scope) {
				break
			}
			bodyNode = bodyNode.Children[1]
		}
		currentValue += 1
	}

	return nil
}

func RunFunc(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	funcVariableName := node.Children[1].Value.Value
	funcParamNode := node.Children[2]
	(*scope)[funcVariableName] = funcParamNode
	return nil
}

func RunReturn(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	if len(node.Children) == 1 {
		(*scope)["RET"] = nil
	}

	mathNode := node.Children[1]
	value, err := RunMath(mathNode, scope)
	if err != nil {
		return err
	}
	(*scope)["RET"] = value
	return nil
}

func RunFuncCall(node *util.TreeNode[parser.ParseNode], scope *Scope) (interface{}, error) {
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
		errorText := fmt.Sprintf(
			"Unmatching number of function paramaters. %s expected %d, got %d.\n",
			funcVariableName, len(funcParamNames), len(mathNodes),
		)
		util.FatalError(errorText, node.Value.Token.Line, node.Value.Token.Column)
	}

	scopeCopy := CopyScope(scope)

	for i := 0; i < len(mathNodes); i++ {
		value, err := RunMath(mathNodes[i], scope)
		if err != nil {
			return nil, err
		}
		varName := funcParamNames[i]
		(*scopeCopy)[varName] = value
	}

	for funcBodyNode.Children[0].Value.Name != "funcend" {
		expressionNode := funcBodyNode.Children[0]
		err := RunExpression(expressionNode, scopeCopy)
		if err != nil {
			return nil, err
		}
		if returnValue, ok := (*scopeCopy)["RET"]; ok {
			return returnValue, nil
		}
		funcBodyNode = funcBodyNode.Children[1]
	}

	return nil, nil
}

func RunPrint(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	result, err := RunMath(node.Children[1], scope)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func RunPrintln(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	result, err := RunMath(node.Children[1], scope)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func RunListAppend(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	if varValue, ok := (*scope)[node.Children[0].Value.Value]; ok {
		varList := varValue.([]interface{})
		newValue, err := RunValue(node.Children[2], scope)
		if err != nil {
			return err
		}
		(*scope)[node.Children[0].Value.Value] = append(varList, newValue)
	} else {
		errorText := fmt.Sprintf("Invalid list variable %s", node.Children[0].Value.Name)
		return util.FormatError(errorText, node.Value.Token.Line, node.Value.Token.Column)
	}
	return nil
}

func RunListPop(node *util.TreeNode[parser.ParseNode], scope *Scope) (interface{}, error) {
	if varValue, ok := (*scope)[node.Children[0].Value.Value]; ok {
		varList := varValue.([]interface{})
		val, err := RunValue(node.Children[2], scope)
		if err != nil {
			return nil, err
		}
		indexValue, ok := val.(int64)
		if !ok {
			util.FatalError("Invalid list index value type", node.Value.Token.Line, node.Value.Token.Column)
		}
		returnValue := varList[indexValue]

		varList = append(varList[:indexValue], varList[indexValue+1:]...)
		(*scope)[node.Children[0].Value.Value] = varList

		return returnValue, nil
	}

	return nil, util.FormatError("Failed to pop value from list", node.Value.Token.Line, node.Value.Token.Column)
}

func RunExpression(node *util.TreeNode[parser.ParseNode], scope *Scope) error {
	for _, childNode := range node.Children {
		switch childNode.Value.Name {
		case "ASSIGN":
			return RunAssign(childNode, scope)
		case "PRINT":
			return RunPrint(childNode, scope)
		case "PRINTLN":
			return RunPrintln(childNode, scope)
		case "IF":
			return RunIf(childNode, scope)
		case "LOOP":
			return RunLoop(childNode, scope)
		case "FUNC":
			return RunFunc(childNode, scope)
		case "FUNCRETURN":
			return RunReturn(childNode, scope)
		case "FUNCCALL":
			_, err := RunFuncCall(childNode, scope)
			return err
		case "LISTAPPEND":
			return RunListAppend(childNode, scope)
		case "LISTPOP":
			_, err := RunListPop(childNode, scope)
			return err
		}
	}

	return nil
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
			err := RunExpression(expressionNode, scope)
			if err != nil {
				util.FatalError(fmt.Sprint(err), expressionNode.Value.Token.Line, expressionNode.Value.Token.Column)
			}
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
