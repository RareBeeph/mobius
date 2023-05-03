package types

import "github.com/faiface/pixel"

type RgbCurve struct {
	ControlPoints []pixel.RGBA
}

func (c *RgbCurve) EvenLagrangeInterp(samplingProgress float64) (sum pixel.RGBA) {
	// When the parameter is 0, the function returns the first control point;
	// when it's 1, the function returns the last control point;
	// when it's 0.5, the function returns the middle control point, and so on.

	// If there is no control point corresponding to the value,
	// the function returns a point along the simplest curve that passes through the control points given the above conditions.

	// The above is equivalent to saying that samplingProgress*(len(c.ControlPoints)-1)
	// can be thought of as an interpolated version of an index into the control points list.

	// Technical considerations:
	// Assume len == 2^n + 1 (2, 3, 5, 9, etc.), since those are the only lengths that can be evenly spaced in our test scheme
	// Consider trying this out with only 4 points. Will it just be stretched in parameter space?
	// Beware Runge's phenomenon for len > 5. Consider Chebyshev spacing or spline interpolation
	// Consider finding a package that implements this more generally or more performantly

	mul := float64(1)
	degree := len(c.ControlPoints) - 1

	// For each control color, find its Lagrange basis polynomial and add the corresponding amount to the sum
	for colIdx, col := range c.ControlPoints {
		mul = 1

		// Construct the basis polynomial by multiplying terms with zeros at all zeroIdx/degree other than colIdx/degree
		for zeroIdx := range c.ControlPoints {
			if colIdx-zeroIdx != 0 {
				mul *= (samplingProgress*float64(degree) - float64(zeroIdx)) / (float64(colIdx - zeroIdx))
			}
		}

		sum = sum.Add(col.Scaled(mul))
	}

	return sum
}
