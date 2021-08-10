package collatz

// Node defines a node in a graph linking numbers in a collatz series
type Node struct {
	EvenParent *Node
	OddParent  *Node

	Value int
}
