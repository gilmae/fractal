package fractal

type Point struct {
	C       complex128
	X       int
	Y       int
	Escape  int
	FinalZ  complex128
	Escaped bool
}
