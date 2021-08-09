package collatz

import (
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

var memo map[int]*Node = make(map[int]*Node)

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

const TiltAngle = 8.0
const LineLength = 60

func DrawLine(x float64, y float64, angle float64, node *Node, parentValue int, dc *gg.Context, stroke bool) {
	if node.Value%2 == 0 {
		angle += TiltAngle * math.Log(2)
	} else {
		angle -= TiltAngle * math.Log(3) * 1.12
	}

	ll := (float64(node.Value) / (1 + math.Pow(float64(node.Value), 1.2))) * LineLength

	rad := angle * (math.Pi / 180)
	x2 := x + ll*math.Cos(rad)
	y2 := y + ll*math.Sin(rad)*-1

	dc.DrawLine(x, y, x2, y2)

	if node.EvenParent != nil {
		DrawLine(x2, y2, angle, node.EvenParent, node.Value, dc, stroke)
	}
	if node.OddParent != nil {
		DrawLine(x2, y2, angle, node.OddParent, node.Value, dc, stroke)
	}
	if stroke && node.EvenParent != nil && node.OddParent != nil {
		dc.SetRGBA(1, 0, 0, 0.1)
		dc.SetLineWidth(6)
		dc.Stroke()
	}
}

func DrawFirstN(n int, filename string) {
	var graph *Node
	for i := 1; i <= n; i++ {
		chain := []int{}

		Calc(i, &chain)
		graph = UpdateGraph(chain, graph)
	}

	const W = 4000
	const H = 2000
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	DrawLine(0, 0.9*H, 0, graph, graph.Value, dc, true)
	dc.SavePNG(filename)
}

func DrawRandomN(n int, max int32, filename string) {
	var graph *Node
	for i := 1; i <= n; i++ {
		chain := []int{}

		Calc(int(rand.Int31n(max)), &chain)
		graph = UpdateGraph(chain, graph)
	}

	const W = 4000
	const H = 2000
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	DrawLine(0, 0.9*H, 0, graph, graph.Value, dc, true)
	dc.SavePNG(filename)
}

func DrawRandomNFullLine(n int, max int32, filename string) {
	const W = 4000
	const H = 2000
	dc := gg.NewContext(W, H)
	dc.SetRGB(rgb(5, 31, 61))
	dc.Clear()

	for i := 1; i <= n; i++ {
		var graph *Node
		chain := []int{}

		Calc(int(rand.Int31n(max)), &chain)
		graph = UpdateGraph(chain, graph)

		dc.SetRGBA(rgba(247, 237, 230, 0.2))

		dc.SetLineWidth(3)
		DrawLine(W*0.05, H*0.8, 0, graph, graph.Value, dc, false)
		dc.Stroke()
	}

	dc.SavePNG(filename)
}

func rgb(r, g, b int) (float64, float64, float64) {
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

func rgba(r, g, b int, a float64) (float64, float64, float64, float64) {
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0, a
}
