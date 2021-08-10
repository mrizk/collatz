package collatz

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

// Params defines the parameters used to draw a vizualization
type Params struct {
	N   int
	Max int32

	Filename string
	Width    int
	Height   int

	TiltAngle             float64
	TiltAngleEvenAdjuster float64
	TiltAngleOddAdjuster  float64

	LineLengthMultiplier float64
	LineWidth            float64
	LineShorteningPower  float64

	BackgroundColor color.Color
	LineColor       color.Color

	StartX     float64
	StartY     float64
	StartAngle float64
}

// Draw creates a visualization with the given params
func Draw(p *Params) {
	dc := gg.NewContext(p.Width, p.Height)
	dc.SetRGBA(rgba(p.BackgroundColor.RGBA()))
	dc.Clear()

	for i := 1; i <= p.N; i++ {
		var graph *Node
		chain := []int{}

		Calc(int(rand.Int31n(p.Max)), &chain)
		graph = UpdateGraph(chain, graph)

		dc.SetRGBA(rgba(p.LineColor.RGBA()))
		dc.SetLineWidth(float64(p.LineWidth))

		DrawLine(p.StartX, p.StartY, 0, graph, dc, p)
		dc.Stroke()
	}

	dc.SavePNG(p.Filename)
}

// DrawLine recursively draws a line from a given starting point in a graph
func DrawLine(x float64, y float64, angle float64, node *Node, dc *gg.Context, p *Params) {
	if node.Value%2 == 0 {
		angle += p.TiltAngle * p.TiltAngleEvenAdjuster
	} else {
		angle -= p.TiltAngle * p.TiltAngleOddAdjuster
	}

	lineLength := (float64(node.Value) / (1 + math.Pow(float64(node.Value), p.LineShorteningPower))) * p.LineLengthMultiplier

	rad := angle * (math.Pi / 180)
	x2 := x + lineLength*math.Cos(rad)
	y2 := y + lineLength*math.Sin(rad)*-1

	dc.DrawLine(x, y, x2, y2)

	if node.EvenParent != nil {
		DrawLine(x2, y2, angle, node.EvenParent, dc, p)
	}
	if node.OddParent != nil {
		DrawLine(x2, y2, angle, node.OddParent, dc, p)
	}
}

func rgba(r, g, b, a uint32) (float64, float64, float64, float64) {
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0, float64(a) / 255.0
}
