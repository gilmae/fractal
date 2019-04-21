package fractal

import (
	mbig "math/big"
	"sync"

	"github.com/gilmae/rescale/big"
	//"sort"
)

func calculate_escape_for_point(p Point, calculator EscapeCalculator) CalculatedPoint {

	iteration, _, escaped := calculator(p.C)

	return CalculatedPoint{p.X, p.Y, iteration, escaped}
}

func plot(base Base, midX float64, midY float64, zoom float64, width int, height int, calculator EscapeCalculator, calculated chan CalculatedPoint) {
	points := make(chan Point, 64)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for p := range points {
				calculated <- calculate_escape_for_point(p, calculator)
			}
			wg.Done()
		}()
	}

	// Derive new bounds based on focal point and zoom
	new_r_start, new_r_end := big.GetZoomedBounds(&base.RMin, &base.RMax, new(mbig.Float).SetFloat64(midX), zoom)
	new_i_start, new_i_end := big.GetZoomedBounds(&base.IMin, &base.IMax, new(mbig.Float).SetFloat64(midY), zoom)

	// Pregenerate all the values of the x  & Y CoOrdinates
	xCoOrds := make([]mbig.Float, width)
	for i, _ := range xCoOrds {
		xCoOrds[i] = big.Rescale(new_r_start, new_r_end, width, i)
	}

	yCoOrds := make([]mbig.Float, height)
	for i, _ := range yCoOrds {
		yCoOrds[height-i-1] = big.Rescale(new_i_start, new_i_end, height, i)
	}

	for x := 0; x < width; x += 1 {
		for y := height - 1; y >= 0; y -= 1 {
			points <- Point{BigComplex{xCoOrds[x], yCoOrds[y]}, x, y, 0, BigComplex{}, false}
		}
	}

	close(points)

	wg.Wait()
}

func Escape_Time_Calculator(base Base, midX float64, midY float64, zoom float64, width int, height int, calculator EscapeCalculator) map[Key]CalculatedPoint {
	var points_map = make(map[Key]CalculatedPoint)

	calculatedChan := make(chan CalculatedPoint)

	go func(points <-chan CalculatedPoint, hash map[Key]CalculatedPoint) {
		for p := range points {
			hash[Key{p.X, p.Y}] = p
		}
	}(calculatedChan, points_map)

	plot(base, midX, midY, zoom, width, height, calculator, calculatedChan)

	return points_map
}

type EscapeCalculator func(z BigComplex) (float64, BigComplex, bool)
