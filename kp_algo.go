package strippacking

import (
	"math/cmplx"
	"container/heap"
)

type Bin struct {
	Rect
	top float64
	t   int
}

type HeapElemLike interface {
	VacantSpace() float64
}

type HeapElem struct {
	p *Bin
}

func (v *HeapElem) VacantSpace() float64 {
	return v.p.H - v.p.top
}

func (v *HeapElem) Less(x interface{}) bool {
	return v.VacantSpace() > x.(HeapElemLike).VacantSpace()
}

// An BinHeap is a min-heap of ints.
type PackingHeap []*HeapElem

func (h PackingHeap) Len() int           { return len(h) }
func (h PackingHeap) Less(i, j int) bool { return h[i].Less(h[j]) }
func (h PackingHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PackingHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*HeapElem))
}

func (h *PackingHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}










// Packs rectangles from slice 'rects' to the strip starting from lower bound 
// 'be' according to Kuzyurin-Pospelov's basic algorithm. 
// Returns an upper bound of resulting alignment. This algo in not quite 
// on-line, it uses number of all the rectangles.
// Ignores 'm' parameter since it is single strip packing algorithm.
type Kp1Algo struct {
	frame    Bin
	bins     map[int]*PackingHeap
	delta, u float64
	d        int
}

func (v *Kp1Algo) Init(n int) {
	v.delta = real(cmplx.Pow(complex(float64(n), 0), (-1.0 / 3)))
	v.u = real(cmplx.Pow(complex(float64(n), 0), (1.0 / 3)))
	v.d = int(1 / (2 * v.delta))
	v.bins = make(map[int]*PackingHeap)
	for y := 0; y <= 2*v.d+1; y++ {
		vec := make(PackingHeap, 0)
		v.bins[y] = &vec
	}
}

func (v *Kp1Algo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	v.frame.top = 0
	v.frame.Y = ybe
	v.frame.X = xbe
	v.frame.W = 1

	n := len(rects)
	v.Init(n)

	for i := 0; i < n; i++ {
		r := &rects[i]
		if r.W > (1 - v.delta) {
			PackToBin(&v.frame, r)
			continue
		}
		
		j := v.RectType(r)
		if v.PackToTopBin(r, j) {
			continue
		}
		// Opening pair of new bins and packing current rectangle into corresponging
		// one.
		v.PackToNewShelfInFrame(r, &v.frame, j)
	}
	return v.frame.Y + v.frame.top
}

// Returns true/false whether rectangle was packed into top-of-the-heap bin.
func (v *Kp1Algo) PackToTopBin(r *Rect, j int) bool {
	if 0 == len(*v.bins[j]) {
		return false
	}
	he := heap.Pop(v.bins[j]).(*HeapElem)
	defer heap.Push(v.bins[j], he)
	if he.VacantSpace() >= r.H {
		PackToBin(he.p, r)
		return true
	}
	return false
}

func (v *Kp1Algo) RectType(r *Rect) int {
	j := int(0)
	for y := 1; y <= v.d; y++ {
		if r.W <= (v.delta * float64(y)) {
			j = y
			break
		}
	}
	if 0 == j {
		for y := v.d; y >= 1; y-- {
			if r.W <= (1 - v.delta*float64(y)) {
				j = v.d*2 - y + 1
				break
			}
		}
	}
	return j
}

func (v *Kp1Algo) PackToNewShelfInFrame(r *Rect, f *Bin, j int) {
	b1 := v.AddBin(j)
	PackToBin(f, &b1.Rect)
	PackToBin(b1, r)
	b2 := v.AddBin(v.ComplType(j))
	b2.Y = b1.Y
	b2.X = b1.X + b1.W
}

func (v *Kp1Algo) AddBin(t int) *Bin {
	b := new(Bin)
	b.H = v.u
	b.W = v.WidthType(t)
	b.top = 0
	b.t = t
	heap.Push(v.bins[t], &HeapElem{b})
	return b
}

func PackToBin(bin *Bin, r *Rect) {
	r.X = bin.X
	r.Y = bin.Y + bin.top
	bin.top += r.H
}

func (v *Kp1Algo) ComplType(t int) int {
	return 2*v.d - t + 1
}

func (v *Kp1Algo) WidthType(t int) float64 {
	if t <= v.d {
		return v.delta * float64(t)
	}
	return 1 - float64(v.ComplType(t))*v.delta
}

// Algorithm on top of Kp1Algo. Kp2 applies Kp1 for consequtive subsequences of 
// length 2, 4, 8, 16, etc in on-line manner so that Kp2Algo actually does not 
// know  total number of rectangles.
type Kp2Algo struct{}

func (v *Kp2Algo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	var b []Rect
	a := rects[:]
	exit_flag := false
	var H float64 = ybe

	for cnt := 2; !exit_flag; cnt *= 2 {
		if cnt > len(a) {
			exit_flag = true
			b = a
		} else {
			b = a[:cnt]
			a = a[cnt:]
		}
		kp1 := new(Kp1Algo)
		H = kp1.Pack(b, xbe, H, m)
	}
	return H
}
