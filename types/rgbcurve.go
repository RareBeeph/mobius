package types

import "github.com/faiface/pixel"

type RgbCurve struct {
	ControlPoints []pixel.RGBA
}

func (c *RgbCurve) EvenLagrangeInterp(t float64) (sum pixel.RGBA) {
	// Assume len == 2^n + 1 (2, 3, 5, 9, etc.), since those are the only lengths that can be evenly spaced in our test scheme
	// Consider trying this out with only 4 points. Will it just be stretched in parameter space?
	// Beware Runge's phenomenon for len > 5. Consider Chebyshev spacing (although that might be hard given our test scheme)
	// Consider finding a package that implements this more generally or more performantly

	mul := float64(1)
	degree := len(c.ControlPoints) - 1

	// For each control color, find its Lagrange basis polynomial and add the corresponding amount to the sum
	for colIdx, col := range c.ControlPoints {
		mul = 1

		// Construct the basis polynomial by multiplying terms with zeros at all zeroIdx/degree other than colIdx/degree
		for zeroIdx := range c.ControlPoints {
			if colIdx-zeroIdx != 0 {
				mul *= (t*float64(degree) - float64(zeroIdx)) / (float64(colIdx - zeroIdx))
			}
		}

		sum = sum.Add(col.Scaled(mul))
	}

	return sum
}
