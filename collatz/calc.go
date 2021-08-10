package collatz

var memo map[int]*Node = make(map[int]*Node)

// Calc caltulates the `chain` of numbers from the given value `n` until 1
func Calc(n int, chain *[]int) {
	*chain = append(*chain, n)

	// if _, ok := memo[n]; ok {
	// 	return
	// }
	memo[n] = nil

	if n == 1 {
		return
	}

	if n%2 == 0 {
		Calc(n/2, chain)
	} else {
		Calc(3*n+1, chain)
	}
}

// UpdateGraph updates the given `graph` with the values in `chain`
func UpdateGraph(chain []int, graph *Node) *Node {
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
