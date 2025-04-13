package util

import "fmt"

type TreeNode[T any] struct {
	Children []*TreeNode[T]
	Value    T
}

var LogLevel int = 0

func Log(verbosity int, text string) {
	if verbosity <= LogLevel {
		fmt.Println(text)
	}
}
