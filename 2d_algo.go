package strippacking

import (
	"math/cmplx"
)

type TdAlgo struct {
	frame Bin
	delta float64
	Rects []Rect
	nrects int
}

func (v *TdAlgo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	v.frame.top = 0
	v.frame.Y = ybe
	v.frame.X = xbe
	v.frame.W = 1 //Confirm
	v.nrects = 0
	n := len(rects)
	//c := complex(float64(n), 0)
	c := complex(float64(n), 0)
	v.delta = real(cmplx.Pow(c, (-1.0 / 2)))
	v.Rects = make([]Rect, n)
	
	for y := 0; y < n; y++ {
		r := &rects[y]
		var best_r *Rect = nil
		var best_s float64 = 0
		var best_vertical bool
		for j := 0; j < v.nrects; j++ {
			packable, vertical := v.Packable(&v.Rects[j], r)
			if packable {
				if (nil != best_r) && (best_s <= v.Rects[j].Area()) {
					continue
				}
				best_r = &v.Rects[j]
				best_s = best_r.Area()
				best_vertical = vertical
			}
		}
		if nil == best_r {
			for j := 0; j < v.nrects; j++ {
				packable, vertical := v.SimplePackable(&v.Rects[j], r)
				if packable {
					//println("splitting")
					v.SplittingPackToRect(&v.Rects[j], r, vertical)
					continue
				}
			}
			println("on top")
			v.PackRectOnTop(r)
		} else {
			println("sane")
			PackToRect(best_r, r, best_vertical)
		}
	}
	println(v.nrects)
	return v.frame.top + v.frame.Y
}

func (v *TdAlgo) SplittingPackToRect(outer *Rect, inner *Rect, vertical bool) {
	inner.X = outer.X
	inner.Y = outer.Y
	side_r := v.AddRect()
	if vertical {
		side_r.X = outer.X
		side_r.Y = outer.Y + inner.H
		side_r.W = inner.W
		side_r.H = outer.H - inner.H
		outer.W -= inner.W
		outer.X += inner.W
	} else {
		// Horizontal.
		side_r.X = outer.X + inner.W
		side_r.Y = outer.Y
		side_r.W = outer.W - inner.W
		side_r.H = inner.H
		outer.H -= inner.H
		outer.Y += inner.H
	}
}

func (v *TdAlgo) AddRect() *Rect {
	v.nrects++
	return &v.Rects[v.nrects - 1]
}

func (v *TdAlgo) PackRectOnTop(r *Rect) {
	outer := v.AddRect()
	outer.W = 1
	outer.H = r.H
	PackToBin(&v.frame, outer)
	PackToRect(outer, r, false)
}

// Returns whether inner fits for packing into outer and corresponding vertical 
// flag. Note that inner is packable when at least one of its dimensions differ 
// from corresponding dimension of outer for no more than delta.
func (v *TdAlgo) Packable(outer *Rect, inner *Rect) (packable bool, vertical bool) {
	if ((inner.W + v.delta) >= outer.W) && (inner.W <= outer.W) && (inner.H <= outer.H) {
		return true, true
	}
	if ((inner.H + v.delta) >= outer.H) && (inner.W <= outer.W) && (inner.H <= outer.H) {
		return true, false
	}
	return false, false
}

func PackToRect(outer *Rect, inner *Rect, vertical bool) {
	if vertical {
		inner.X = outer.X
		inner.Y = outer.Y
		outer.Y += inner.H
		outer.H -= inner.H
	} else {
		// Horizontal then.
		inner.X = outer.X
		inner.Y = outer.Y
		outer.X += inner.W
		outer.W -= inner.W
	}
}

func (v *Rect) Area() float64 {
	return v.X * v.Y
}

func (v *TdAlgo) SimplePackable(outer *Rect, inner *Rect) (packable bool, vertical bool) {
	if (inner.W > outer.W) || (inner.H > outer.H) {
		return false, false
	}
	if (outer.H - inner.H) > (outer.W - inner.W) {
		return true, false
	}
	return true, true
}