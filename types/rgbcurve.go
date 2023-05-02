package types

import "github.com/faiface/pixel"

type RgbCurve struct {
	ControlPoints []pixel.RGBA
}

func (c *RgbCurve) EvenLagrangeInterp(t float64) (sum pixel.RGBA) {
	/*
		Brainstorming

		if len < 2, something is wrong
		if len == 2, f = (1-t)*[0] + t*[1]
		if len == 3,
			f(0) = [0]
			f(0.5) = [1]
			f(1) = [2]

			f = at^2 + bt + c
			c = [0]

			0.25a + 0.5b + c = [1]
			b = -3*[0] + 4*[1] - [2]

			a = 2*[0] - 4*[1] + 2*[2]

			so f = (2t^2 - 3t + 1)*[0] + (-4t^2 + 4t)*[1] + (2t^2 - t)*[2]


			((x-1/2)/(-1/2))((x-1)/(-1)) = 2x^2 - 3x + 1
			plug in t for x, multiply by [i], and sum


			OR

			construct an n+1 by n+1 matrix based on plugging in i/n for t

			|  1   1  1| |a|   |[0]|
			|0.25 0.5 1| |b| = |[1]|
			|  0   0  1| |c|   |[2]|

			find the inverse matrix (there's probably plenty of great implementations but i can't be bothered to find them)

			     |[0]|	 |a|
			M^-1 |[1]| = |b|
			     |[2]|	 |c|
	*/

	// Assume len == 2^n + 1 (2, 3, 5, 9, etc.), since those are the only lengths that can be evenly spaced in our test scheme
	// Beware Runge's phenomenon for len > 5. Consider Chebyshev spacing (although that might be hard given our test scheme)
	// Consider finding a package that implements this more generally or more performantly

	mul := float64(1)
	n := len(c.ControlPoints) - 1

	// For each control color, find its Lagrange basis polynomial and add the corresponding amount to the sum
	for i, col := range c.ControlPoints {
		mul = 1

		// Construct the basis polynomial by multiplying terms with zeros at all j/n other than i/n
		for j := range c.ControlPoints {
			if i-j != 0 {
				mul *= (t - float64(j)/float64(n)) / (float64(i-j) / float64(n))
			}
		}

		sum = sum.Add(col.Scaled(mul))
	}

	return sum
}
