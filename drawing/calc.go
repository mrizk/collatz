package drawing

// Node defines a node in a graph linking numbers in a collatz series
type Node struct {
	EvenParent *Node
	OddParent  *Node

	Value int
}

// Calc caltulates the `chain` of numbers from the given value `n` until 1
func Calc(n int, chain *[]int, memo map[int]*Node) {
	*chain = append(*chain, n)

	// if _, ok := memo[n]; ok {
	// 	return
	// }
	memo[n] = nil

	if n == 1 {
		return
	}

	if n%2 == 0 {
		Calc(n/2, chain, memo)
	} else {
		Calc(3*n+1, chain, memo)
	}
}

// UpdateGraph updates the given `graph` with the values in `chain`
func UpdateGraph(chain []int, graph *Node, memo map[int]*Node) *Node {
	if graph == nil {
		graph = &Node{
			Value: 1,
		}
	}
	currentNode := graph
	for i := len(chain) - 1; i >= 0; i-- {
		if node, ok := memo[chain[i]]; ok && node != nil {
			currentNode = memo[chain[i]]
		} else {
			if chain[i]%2 == 0 {
				if currentNode.EvenParent == nil {
					currentNode.EvenParent = &Node{
						Value: chain[i],
					}
				}
				currentNode = currentNode.EvenParent
			} else if chain[i] != 1 {
				if currentNode.OddParent == nil {
					currentNode.OddParent = &Node{
						Value: chain[i],
					}
				}
				currentNode = currentNode.OddParent
			}
		}
		memo[currentNode.Value] = currentNode
	}
	return graph
}
