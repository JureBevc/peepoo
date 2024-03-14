package util

import "fmt"

type TreeNode[T any] struct {
	Children []*TreeNode[T]
	Value    T
}

func PrintTree[T any](tree *TreeNode[T], prefix string) {
	fmt.Printf("%s%s\n", prefix, fmt.Sprint(tree.Value))
	childPrefix := prefix + prefix
	for _, child := range (*tree).Children {
		PrintTree(child, childPrefix)
	}
}
