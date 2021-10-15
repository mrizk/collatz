package api_test

import (
	"collatz/api"
	"image/color"
	"math"
	"testing"
)

func TestDraw(t *testing.T) {
	width := 15118
	height := 8504
	p := &api.Params{
		N:   3000,
		Max: 200000000,

		Filename: "test",
		Format:   api.ImageFormatPNG,
		Width:    width,
		Height:   height,

		TiltAngle:             8.0,
		TiltAngleEvenAdjuster: math.Log(2),
		TiltAngleOddAdjuster:  math.Log(3) * 1.12,

		LineLengthMultiplier: 260,
		LineWidth:            7,
		LineShorteningPower:  1.21,

		BackgroundColor: color.RGBA{0, 0, 0, 0},
		LineColor:       color.RGBA{0, 0, 0, 60},
		// BackgroundColor: color.RGBA{255, 249, 245, 255},
		// LineColor:       color.RGBA{255, 230, 191, 60},
		// GradientColors: []color.Color{
		// 	color.RGBA{12, 36, 77, 255}, // red 235, 64, 52
		// 	color.RGBA{0, 0, 0, 255},
		// },

		GridColor:     color.RGBA{7, 67, 110, 255},
		GridLineWidth: 4,
		// GridColumns:   64,
		// GridRows:      36,
		GridColumns: 0,
		GridRows:    0,

		StartX: float64(width) * 0.0,
		StartY: float64(height) * 0.80,
	}

	// data, err := json.MarshalIndent(p, "", "	")
	// if err != nil {
	// 	panic(err)
	// }
	// println(string(data))

	api.Draw(p)
}
