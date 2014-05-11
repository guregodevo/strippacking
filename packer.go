package strippacking

import (
	"math/rand"
	"fmt"
	"math/cmplx"
)

type Rect struct {
	X, Y float64
	H, W float64
}

func GenerateRectangles(n int) []Rect {
	res := make([]Rect, n)
	for y := 0; y < n; y++ {
		res[y].H = rand.Float64()
		res[y].W = rand.Float64()
		res[y].X = rand.Float64()
		res[y].Y = rand.Float64()
	}
	return res
}

func TotalArea(rects []Rect) float64 {
	var res float64 = 0
	for _, r := range rects {
		res += r.W * r.H
	}
	return res
}

type Algorithm interface {
	Pack(rects []Rect, xbe, ybe float64, m int) float64
}

// Returns uncovered area divided by N^(2/3).
func Run(n int, render, validate bool, algo_name string, m int) (coefficient float64) {
	rects := GenerateRectangles(n)
	var algo Algorithm

	if "kp1" == algo_name {
		algo = new(Kp1Algo)
	} else if "kp2" == algo_name {
		algo = new(Kp2Algo)
	} else if "2d" == algo_name {
		algo = new(TdAlgo)
	} else if "kp2_msp" == algo_name {
		algo = new(Kp2MspAlgo)
	} else if "kp2_msp_b" == algo_name {
		algo = new(Kp2MspBalancedAlgo)
	} else {
		algo = new(Kp2Algo)
	}

	H = algo.Pack(rects, 0, 0, m)
	total_area := TotalArea(rects)
	fmt.Printf("Solution height = %0.9v\nTotal area = %0.9v\n", H, total_area)
	uncovered_area := H * float64(m) - total_area
	fmt.Printf("Uncovered area = %0.9v\n", uncovered_area)
	coefficient = uncovered_area / real(cmplx.Pow(complex(float64(n), 0), (2.0/3)))

	if true == validate {
		if false == Validate(rects, m, H) {
			println("Validation: ERROR")
		} else {
			println("Validation: OK")
		}
	}

	if render {
		render_all(rects, m)
	}
	return
}

