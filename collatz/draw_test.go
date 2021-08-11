package collatz_test

import (
	"collatz/collatz"
	"image/color"
	"math"
	"testing"
)

func TestDraw(t *testing.T) {
	width := 14173
	height := 8031
	p := &collatz.Params{
		N:   3000,
		Max: 200000000,

		Filename: "out.png",
		Width:    width,
		Height:   height,

		TiltAngle:             8.0,
		TiltAngleEvenAdjuster: math.Log(2),
		TiltAngleOddAdjuster:  math.Log(3) * 1.12,

		LineLengthMultiplier: 250,
		LineWidth:            7,
		LineShorteningPower:  1.21,

		BackgroundColor: color.RGBA{255, 249, 245, 255},
		LineColor:       color.RGBA{18, 69, 128, 60},

		StartX: float64(width) * 0.0,
		StartY: float64(height) * 0.85,
	}

	collatz.Draw(p)
}
