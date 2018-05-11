package fractal

import (
	"sync"
    "github.com/gilmae/rescale"
    //"sort"
)

  
func calculate_escape_for_point(p Point, calculator EscapeCalculator) Point {
	
	iteration, finalZ, escaped := calculator(p.C);
	
	return Point{p.C, p.X, p.Y, iteration, finalZ, escaped}
}
  
func plot(base Base, midX float64, midY float64, zoom float64, width int, height int, calculator EscapeCalculator, calculated chan Point) {
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
	new_r_start, new_r_end := rescale.GetZoomedBounds(base.RMin, base.RMax, midX, zoom)
	new_i_start, new_i_end := rescale.GetZoomedBounds(base.IMin, base.IMax, midY, zoom)
  
	// Pregenerate all the values of the x  & Y CoOrdinates
	xCoOrds := make([]float64, width)
	for i,_ := range xCoOrds {
	  xCoOrds[i] = rescale.Rescale(new_r_start, new_r_end, width, i);
	}
  
	yCoOrds := make([]float64, height)
	for i,_ := range yCoOrds {
	  yCoOrds[height-i-1] = rescale.Rescale(new_i_start, new_i_end, height, i);
	}
  
	for x:=0; x < width; x += 1 {
	  for y:=height-1; y >= 0; y -= 1 {
		points <- Point{complex(xCoOrds[x], yCoOrds[y]),x,y, 0, complex(0,0), false}
	  }
	}
  
	close(points)
  
	wg.Wait()
  }

  func Escape_Time_Calculator(base Base, midX float64, midY float64, zoom float64, width int, height int, calculator EscapeCalculator) map[Key]Point {
	var points_map = make(map[Key]Point)
  
	calculatedChan := make(chan Point)
  
	go func(points<-chan Point, hash map[Key]Point) {
	  for p := range points {
		hash[Key{p.X,p.Y}] = p
	  }
	}(calculatedChan, points_map)
  
	plot(base, midX, midY, zoom, width, height, calculator, calculatedChan)
	  
	return points_map
  }

  type EscapeCalculator func(z complex128) (float64, complex128, bool)

