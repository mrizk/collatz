package collatz_test

import (
	"collatz/collatz"
	"testing"
)

func TestDraw(t *testing.T) {
	// collatz.DrawFirstN(5000, "out.png")

	// collatz.DrawRandomN(1000, 20000000, "out.png")

	collatz.DrawRandomNFullLine(1500, 200000000, "out.png")
}
