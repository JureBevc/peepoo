package util

type TreeNode[T any] struct {
	Children []*TreeNode[T]
	Value    T
}
