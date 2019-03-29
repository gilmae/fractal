package fractal

type Point struct {
	C       BigComplex
	X       int
	Y       int
	Escape  float64
	FinalZ  BigComplex
	Escaped bool
}

type CalculatedPoint struct {
	X       int
	Y       int
	Escape  float64
	Escaped bool
}
