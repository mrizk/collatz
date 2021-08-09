package collatz

type Node struct {
	EvenParent *Node
	OddParent  *Node

	Value int
}
