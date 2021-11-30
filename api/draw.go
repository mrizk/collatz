package api

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/opentype"
)

// Params defines the parameters used to draw a vizualization
type Params struct {
	N   int
	Max int32

	Filename string
	Format   ImageFormat
	Width    int
	Height   int

	TiltAngle             float64
	TiltAngleEvenAdjuster float64
	TiltAngleOddAdjuster  float64

	LineLengthMultiplier float64
	LineWidth            float64
	LineShorteningPower  float64

	BackgroundColor color.RGBA
	LineColor       color.RGBA
	TextColor       color.RGBA
	GradientColors  []color.RGBA

	GridColor     color.RGBA
	GridLineWidth float64
	GridColumns   int
	GridRows      int

	StartX     float64
	StartY     float64
	StartAngle float64
}

type ImageFormat string

const (
	ImageFormatPNG ImageFormat = "png"
	ImageFormatJPG ImageFormat = "jpg"
)

// Draw creates a visualization with the given params
func Draw(p *Params) (image.Image, error) {
	// Init
	dc := gg.NewContext(p.Width, p.Height)
	dc.SetRGBA(rgba(p.BackgroundColor.RGBA()))
	dc.Clear()

	// Gradient
	if len(p.GradientColors) > 0 {
		grad := gg.NewRadialGradient(float64(p.Width)/2.0, float64(p.Height)/2.0, 0, float64(p.Width)/2.0, float64(p.Height)/2.0, float64(p.Width)*2)
		for i, c := range p.GradientColors {
			grad.AddColorStop(float64(i), c)
		}
		dc.SetFillStyle(grad)
		dc.DrawRectangle(0, 0, float64(p.Width), float64(p.Height))
		dc.Fill()
	}

	// Gird
	dc.SetColor(p.GridColor)
	dc.SetLineWidth(p.GridLineWidth)

	columns := p.GridColumns
	for i := 1; i <= columns; i++ {
		width := p.Width / columns
		dc.DrawLine(float64(i*width), 0, float64(i*width), float64(p.Height))
		dc.Stroke()
	}
	rows := p.GridRows
	for i := 1; i <= rows; i++ {
		height := p.Height / rows
		dc.DrawLine(0, float64(i*height), float64(p.Width), float64(i*height))
		dc.Stroke()
	}

	// Lines
	for i := 1; i <= p.N; i++ {
		var graph *Node
		chain := []int{}
		memo := make(map[int]*Node)

		Calc(int(rand.Int31n(p.Max)), &chain, memo)
		graph = UpdateGraph(chain, graph, memo)

		dc.SetRGBA(rgba(p.LineColor.RGBA()))
		dc.SetLineWidth(float64(p.LineWidth))

		DrawLine(p.StartX, p.StartY, 0, graph, dc, p)
		dc.Stroke()
	}

	// Signature
	dc.SetRGBA(rgba(p.TextColor.RGBA()))
	fontBytes, err := ioutil.ReadFile("Magis Authentic.ttf")
	if err != nil {
		return nil, err
	}
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: 100,
		DPI:  150,
		// Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(face)
	dc.DrawString("No -> 3n + 1", float64(p.Width)-1250, float64(p.Height)-500)
	dc.DrawString("Ne -> n/2", float64(p.Width)-1250, float64(p.Height)-270)

	// Save
	filename := fmt.Sprintf("images/%s.%s", p.Filename, string(p.Format))
	switch p.Format {
	case ImageFormatPNG:
		err := dc.SavePNG(filename)
		if err != nil {
			return nil, err
		}
	case ImageFormatJPG:
		err := gg.SaveJPG(filename, dc.Image(), 100)
		if err != nil {
			return nil, err
		}
	}
	return dc.Image(), nil
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
