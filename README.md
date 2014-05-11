HOW-TO
=======

```go
package main

import (
	"github.com/guregodevo/strippacking"
	"math/rand"
	"time"
	"fmt"
	"flag"
	"math/cmplx"
)

func main() {
	prender := flag.Bool("r", false, "Render resulting alignment of all the rectangles")
	//prenderbins := flag.Bool("rb", false, "Render bins")
	//pnonsolid := flag.Bool("ns", false, "Non solid rendering of rectangles")
	pn := flag.Int("n", 100, "Number of rectangles")
	pm := flag.Int("m", 1, "Number of strips")
	pvalidate := flag.Bool("v", false, "Validate resulting alignment")
	palgo := flag.String("a", "kp2", "Type of algorithm")
	ptimes := flag.Int("t", 1, "Number of tests")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	println("Number of rectangles = ", *pn)
	fmt.Printf("N^(2/3) = %0.9v\n\n", real(cmplx.Pow(complex(float64(*pn), 0), (2.0/3))))

	var coef_s float64 = 0
	for y := 0; y < *ptimes; y++ {
		coef := strippacking.Run(*pn, *prender, *pvalidate, *palgo, *pm)
		coef_s += coef
	}
	fmt.Printf("\nAverage coefficient = %0.9v\n", coef_s/float64(*ptimes))
}
```
