package collatz_test

import (
	"collatz/collatz"
	"image/color"
	"math"
	"testing"
)

func TestDraw(t *testing.T) {
	width := 4000
	height := 2000
	p := &collatz.Params{
		N:   1500,
		Max: 200000000,

		Filename: "out.png",
		Width:    width,
		Height:   height,

		TiltAngle:             8.0,
		TiltAngleEvenAdjuster: math.Log(2),
		TiltAngleOddAdjuster:  math.Log(3) * 1.12,

		LineLengthMultiplier: 60,
		LineWidth:            3,
		LineShorteningPower:  1.2,

		BackgroundColor: color.RGBA{5, 31, 61, 255},
		LineColor:       color.RGBA{247, 237, 230, 51},

		StartX: float64(width) * 0.05,
		StartY: float64(height) * 0.8,
	}

	collatz.Draw(p)
}
