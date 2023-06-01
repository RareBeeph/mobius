package types

import (
	"log"
	"math"
)

type slopes struct {
	dudR, dvdR, dudG, dvdG, dudB, dvdB float64
}

func SlopesOfMetric(metric [3][3]float64) (out slopes) {
	/*
		primed basis vectors defined to be tangent to the manifold, with R G and B components maintained from pixel's RGB parametrization
		R' = R^ + dudR*u^ + dvdR*v^
		G' = G^ + dudG*u^ + dvdG*v^
		B' = B^ + dudB*u^ + dvdB*v^

		metric components are dot products of corresponding R' (index 0), G' (index 1), B' (index 2)
		metric[0][0] = 1 + dudR^2 + dvdR^2
		metric[1][1] = 1 + dudG^2 + dvdG^2
		metric[2][2] = 1 + dudB^2 + dvdB^2
		metric[0][1] = dudR*dudG + dvdR*dvdG
		metric[0][2] = dudR*dudB + dvdR*dvdB
		metric[1][2] = dudG*dudB + dvdG*dvdB

		metric components are constant and known, so we can solve the diagonals for dudR in terms of dvdR, etc
		dudR = +-sqrt(metric[0][0] - 1 - dvdR^2)
		dudG = +-sqrt(metric[1][1] - 1 - dvdG^2)
		dudB = +-sqrt(metric[2][2] - 1 - dvdB^2)

		plug those into the off diagonals. this gives us 3 equations for elliptic cylinders along perpendicular axes
		metric[0][1] = +-sqrt((metric[0][0] - 1 - dvdR^2)*(metric[1][1] - 1 - dvdG^2)) + dvdR*dvdG
		metric[0][2] = +-sqrt((metric[0][0] - 1 - dvdR^2)*(metric[2][2] - 1 - dvdB^2)) + dvdR*dvdB
		metric[1][2] = +-sqrt((metric[1][1] - 1 - dvdG^2)(metric[2][2] - 1 - dvdB^2))B + dvdG*dvdB

		rearrange, square, distribute, and combine like terms to kill the sqrt and make the expressions comparable
		metric[0][1]^2 - 2*dvdR*dvdG*metric[0][1] = metric[0][0]*metric[1][1] - metric[0][0] - metric[0][0]*dvdG^2 - metric[1][1] + 1 + dvdG^2 - metric[1][1]*dvdR^2 + dvdR^2
		etc

		rearrange more to make the equations comparable
		0 = (1 - metric[0][0])*dvdG^2 + 2*metric[0][1]*dvdR*dvdG + (1 - metric[1][1])*dvdR^2 + (1 - metric[0][0])*(1 - metric[1][1]) - metric[0][1]^2
		0 = (1 - metric[0][0])*dvdB^2 + 2*metric[0][2]*dvdR*dvdB + (1 - metric[2][2])*dvdR^2 + (1 - metric[0][0])*(1 - metric[2][2]) - metric[0][2]^2
		0 = (1 - metric[1][1])*dvdB^2 + 2*metric[1][2]*dvdG*dvdB + (1 - metric[2][2])*dvdG^2 + (1 - metric[1][1])*(1 - metric[2][2]) - metric[1][2]^2

		solve the latter two ellipse equations for dvdR and dvdG in terms of dvdB--take the positive branch sqrt
		dvdR = -(metric[0][2]/(1 - metric[2][2]))*dvdB + sqrt((metric[0][2]^2 + (metric[0][2]^2/(1 - metric[2][2]) + metric[0][0] - 1)*dvdB^2)/(1 - metric[2][2]) + metric[0][0] - 1)
		dvdG = -(metric[1][2]/(1 - metric[2][2]))*dvdB + sqrt((metric[1][2]^2 + (metric[1][2]^2/(1 - metric[2][2]) + metric[1][1] - 1)*dvdB^2)/(1 - metric[2][2]) + metric[1][1] - 1)

		manipulate those
		dvdR = (metric[0][2]/(metric[2][2]-1))*dvdB + sqrt((metric[0][0]-1)/(metric[2][2]-1) - metric[0][2]^2/(metric[2][2]-1)^2) * sqrt(metric[2][2]-1-dvdB^2)
		dvdG = (metric[1][2]/(metric[2][2]-1))*dvdB + sqrt((metric[1][1]-1)/(metric[2][2]-1) - metric[1][2]^2/(metric[2][2]-1)^2) * sqrt(metric[2][2]-1-dvdB^2)

		plug into the first ellipse equation and solve for dvdB in terms of the metric components--take the positive branch sqrt
	*/
	metric = metricScale(metric, 10)

	// 0 = (1 - metric[0][0])*dvdG^2 + 2*metric[0][1]*dvdR*dvdG + (1 - metric[1][1])*dvdR^2 + (1 - metric[0][0])*(1 - metric[1][1]) - metric[0][1]^2
	c11 := 1 - metric[0][0]
	c12 := 2 * metric[0][1]
	c13 := 1 - metric[1][1]
	c14 := (1-metric[0][0])*(1-metric[1][1]) - math.Pow(metric[0][1], 2)

	// dvdR = (metric[0][2]/(metric[2][2]-1))*dvdB + sqrt((metric[0][0]-1)/(metric[2][2]-1) - metric[0][2]^2/(metric[2][2]-1)^2) * sqrt(metric[2][2]-1-dvdB^2)
	c21 := metric[0][2] / (metric[2][2] - 1)
	c22 := math.Sqrt((metric[0][0]-1)/(metric[2][2]-1) - math.Pow(metric[0][2], 2)/math.Pow(metric[2][2]-1, 2))

	// dvdG = (metric[1][2]/(metric[2][2]-1))*dvdB + sqrt((metric[1][1]-1)/(metric[2][2]-1) - metric[1][2]^2/(metric[2][2]-1)^2) * sqrt(metric[2][2]-1-dvdB^2)
	c31 := metric[1][2] / (metric[2][2] - 1)
	c32 := math.Sqrt((metric[1][1]-1)/(metric[2][2]-1) - math.Pow(metric[1][2], 2)/math.Pow(metric[2][2]-1, 2))

	// shared
	c3 := metric[2][2] - 1

	// plug in dvdG and dvdR into eq. 1--it gets LONG
	c41 := c11*math.Pow(c31, 2) - c11*math.Pow(c32, 2) + c12*c31*c21 - c12*c32*c22 + c13*math.Pow(c21, 2) - c13*math.Pow(c22, 2) // dvdB^2 coeff
	c42 := 2*c11*c31*c32 + c12*c32*c21 + c12*c31*c22 + 2*c13*c22*c21                                                             // dvdB*sqrt(c3-dvdB^2) coeff
	c43 := c12*c32*c22*c3 + c11*math.Pow(c32, 2)*c3 + c13*math.Pow(c22, 2)*c3 + c14                                              // constant

	// rearrange as a quadratic of dvdB^2
	c51 := math.Pow(c41, 2) + math.Pow(c42, 2) // dvdB^4 coeff
	c52 := 2*c43*c41 - c3*math.Pow(c42, 2)     // dvdB^2 coeff
	c53 := math.Pow(c43, 2)                    // constant

	out.dvdB = math.Sqrt((-c52 + math.Sqrt(math.Pow(c52, 2)-4*c51*c53)) / (2 * c51)) // quadratic formula on eq. 5
	out.dvdR = c21*out.dvdB + c22*math.Sqrt(c3-math.Pow(out.dvdB, 2))                // eq.2
	out.dvdG = c31*out.dvdB + c32*math.Sqrt(c3-math.Pow(out.dvdB, 2))                // eq.3
	out.dudR = math.Sqrt(metric[0][0] - 1 - math.Pow(out.dvdR, 2))                   // definition of R'
	out.dudG = math.Sqrt(metric[1][1] - 1 - math.Pow(out.dvdG, 2))                   // definition of G'
	out.dudB = math.Sqrt(metric[2][2] - 1 - math.Pow(out.dvdB, 2))                   // definition of B'

	// debug
	var test [3][3]float64
	test[0][0] = 1 + math.Pow(out.dudR, 2) + math.Pow(out.dvdR, 2)
	test[1][1] = 1 + math.Pow(out.dudG, 2) + math.Pow(out.dvdG, 2)
	test[2][2] = 1 + math.Pow(out.dudB, 2) + math.Pow(out.dvdB, 2)
	test[0][1] = -out.dudR*out.dudG + out.dvdR*out.dvdG
	test[1][0] = test[0][1]
	test[0][2] = -out.dudR*out.dudB + out.dvdR*out.dvdB
	test[2][0] = test[0][2]
	test[1][2] = -out.dudG*out.dudB + out.dvdG*out.dvdB
	test[2][1] = test[1][2]
	log.Println(metric)
	log.Println(test)

	return out
}

func metricScale(metric [3][3]float64, factor float64) (out [3][3]float64) {
	for i := range metric {
		for j := range metric {
			out[i][j] = factor * metric[i][j]
		}
	}
	return out
}
