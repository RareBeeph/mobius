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
	metric = metricScale(metric, 100) // Scale up to ensure a valid unprojection. Shouldn't affect local relationships.

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

	// condition logic time. we gotta pick the right branches of the sqrts.

	// setting up a quadratic equation for the metric[0][2] values at which the outer phase changes occur
	d1 := (metric[1][1] - 1) / (metric[2][2] - 1)                                                                                     // metric[0][2]^2 coeff
	d2 := 2 * metric[0][1] * math.Sqrt((metric[1][1]-1)*(metric[2][2]-1)-math.Pow(metric[1][2], 2)) / (metric[2][2] - 1)              // sqrt(d3 - metric[0][2]^2) coeff
	d3 := (metric[2][2] - 1) * (metric[0][0] - 1)                                                                                     // internal value in the sqrt; acts as a radius^2 term
	d4 := (metric[0][0]-1)*math.Pow(metric[1][2], 2)/(metric[2][2]-1) - (metric[0][0]-1)*(metric[1][1]-1) - math.Pow(metric[0][1], 2) // constant

	log.Printf("d: %f, %f, %f, %f", d1, d2, d3, d4)

	// rearrange and square the above to get a quadratic of metric[0][2]^2
	d21 := math.Pow(d1, 2)                      // metric[0][2]^4 coeff
	d22 := math.Pow(d2, 2) + 2*d1*d4            // metric[0][2]^2 coeff
	d23 := math.Pow(d4, 2) - math.Pow(d2, 2)*d3 // constant

	log.Printf("d2: %f, %f, %f", d21, d22, d23)

	// set up a quadratic equation in terms of metric[1][2] for when the lowest phase change occurs at metric[0][2] == 0
	d31 := math.Pow(metric[0][0]-1, 2) / math.Pow(metric[2][2]-1, 2)                                                                                           // metric[1][2]^4 coeff
	d32 := -2*math.Pow(metric[0][0]-1, 2)*(metric[1][1]-1)/(metric[2][2]-1) + 2*(metric[0][0]-1)*math.Pow(metric[0][1], 2)/(metric[2][2]-1)                    // metric[1][2]^2 coeff
	d33 := math.Pow(metric[0][0]-1, 2)*math.Pow(metric[1][1]-1, 2) + math.Pow(metric[0][1], 4) - 2*math.Pow(metric[0][1], 2)*(metric[0][0]-1)*(metric[1][1]-1) // constant

	log.Printf("d3: %f, %f, %f", d31, d32, d33)

	// quadratic formula of above, compared with metric[1][2]. controls the sign of the sqrt for metric[0][2]
	bouncer := math.Copysign(1, -metric[1][2]*math.Copysign(1, metric[0][1])-math.Sqrt((-d32+math.Sqrt(math.Pow(d32, 2)-4*d31*d33))/(2*d31)))

	log.Printf("bouncer: %f", bouncer)

	// quadratic formula of d21 through d23
	lowerphase := math.Sqrt((-d22+math.Copysign(1, metric[0][1]*metric[1][2])*math.Sqrt(math.Pow(d22, 2)-4*d21*d23))/(2*d21)) * bouncer
	upperphase := -lowerphase

	// simplified from an expression of c42 by exploiting invariance
	middlephase := -metric[1][2] * math.Sqrt((metric[0][0]-1)/(metric[1][1]-1))

	log.Printf("phases: %f, %f, %f", lowerphase, middlephase, upperphase)

	z := math.Sqrt((-c52 + math.Sqrt(math.Pow(c52, 2)-4*c51*c53)) / (2 * c51))  // quadratic formula on eq. 5
	z2 := math.Sqrt((-c52 - math.Sqrt(math.Pow(c52, 2)-4*c51*c53)) / (2 * c51)) // negative branch of inner sqrt
	offset := math.Sqrt(c3 - math.Pow(z, 2))
	offset2 := math.Sqrt(c3 - math.Pow(z2, 2))

	log.Printf("calculations: %f, %f, %f, %f", z, z2, offset, offset2)

	// this expression represents the value of metric[1][2] at which the outer phase changes occur at the maximum magnitude of metric[0][2]
	// metric[1][2] being above it or below it determines which of the lowest or the highest phase changes occurs
	if metric[1][2] <= math.Sqrt(math.Pow(metric[0][1], 2)*(metric[2][2]-1)/(metric[0][0]-1)) {
		if metric[0][2] <= lowerphase {
			log.Println("black")
			out.dvdR = -c21*z2 + c22*offset2 // eq.2
			out.dvdG = -c31*z2 + c32*offset2 // eq.3
			out.dvdB = z2
		} else if metric[0][2] <= middlephase {
			log.Println("blackside red")
			// positive branch of z2 outer sqrt
			out.dvdR = c21*z2 + c22*offset2
			out.dvdG = c31*z2 + c32*offset2
			out.dvdB = z2
		} else {
			log.Println("blackside orange")
			// positive branch of z inner sqrt as well
			out.dvdR = c21*z + c22*offset
			out.dvdG = c31*z + c32*offset
			out.dvdB = z
		}
	} else {
		if metric[0][2] <= middlephase {
			log.Println("greenside red")
			// same equations as the middle section of the metric[1][2] <= case, but without a lowest phase interfering
			out.dvdR = c21*z2 + c22*offset2
			out.dvdG = c31*z2 + c32*offset2
			out.dvdB = z2
		} else if metric[0][2] <= upperphase {
			log.Println("greenside orange")
			// same equations as the final section of the metric[1][2] <= case, but with a highest phase interfering
			out.dvdR = c21*z + c22*offset
			out.dvdG = c31*z + c32*offset
			out.dvdB = z
		} else {
			log.Println("green")
			// negative branch of offset outer sqrt
			out.dvdR = c21*z - c22*offset
			out.dvdG = c31*z - c32*offset
			out.dvdB = z
		}
	}

	// currently uncooperative--these need to have their signs constrained, but sometimes the constraints are unsatisfiable!
	out.dudR = math.Sqrt(metric[0][0] - 1 - math.Pow(out.dvdR, 2)) // definition of R'
	out.dudG = math.Sqrt(metric[1][1] - 1 - math.Pow(out.dvdG, 2)) // definition of G'
	out.dudB = math.Sqrt(metric[2][2] - 1 - math.Pow(out.dvdB, 2)) // definition of B'

	// debug
	var test [3][3]float64
	test[0][0] = 1 + math.Pow(out.dudR, 2) + math.Pow(out.dvdR, 2)
	test[1][1] = 1 + math.Pow(out.dudG, 2) + math.Pow(out.dvdG, 2)
	test[2][2] = 1 + math.Pow(out.dudB, 2) + math.Pow(out.dvdB, 2)
	test[0][1] = out.dudR*out.dudG + out.dvdR*out.dvdG
	test[1][0] = test[0][1]
	test[0][2] = out.dudR*out.dudB + out.dvdR*out.dvdB
	test[2][0] = test[0][2]
	test[1][2] = out.dudG*out.dudB + out.dvdG*out.dvdB
	test[2][1] = test[1][2]
	log.Println(metric)
	log.Println(test)
	log.Printf("%f, %f, %f, %f, %f, %f", out.dudR, out.dudG, out.dudB, out.dvdR, out.dvdG, out.dvdB)

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
