package fractal

import "math/big"

type BigComplex struct {
	R, I big.Float
}

/* Multiplying two complex numbers looks like a quadratic equation
 *
 * (x + yi)(w + zi) == w * w + xzi + wyi + yzi^2
 * And since i^2 = -1, this becomes ww-yz + (xz+wy)i
 */
func (z *BigComplex) Mul(x, y *BigComplex) *BigComplex {
	rSquared := new(big.Float).Copy(&x.R)

	rSquared.Mul(rSquared, &y.R)

	iSquared := new(big.Float).Copy(&x.I)
	iSquared.Mul(iSquared, &y.I)

	r := new(big.Float).Copy(rSquared)
	r.Sub(r, iSquared)

	i := new(big.Float).Copy(&x.R)
	i.Mul(i, &y.I)

	tmp := new(big.Float).Copy(&x.I)
	tmp.Mul(tmp, &y.R)

	i.Add(i, tmp)

	z.R = *r
	z.I = *i

	return z
}

func (z *BigComplex) Add(x, y *BigComplex) *BigComplex {
	z.R.Add(&x.R, &y.R)
	z.I.Add(&x.I, &y.I)

	return z
}

func (z *BigComplex) Sub(x, y *BigComplex) *BigComplex {
	z.R.Sub(&x.R, &y.R)
	z.I.Sub(&x.I, &y.I)

	return z
}

/*
	x + yi / u + vi = (xu + yv) + (-xv + yu)i / u^2 + v^2
*/
func (z *BigComplex) Quo(x, y *BigComplex) *BigComplex {
	divisor := new(big.Float).Set(&y.R)
	divisor.Mul(divisor, divisor) // u^2
	tmp := new(big.Float).Set(&y.I)
	tmp.Mul(tmp, tmp)         // v^2
	divisor.Add(divisor, tmp) // u^2 + v^2

	xu := new(big.Float).SetFloat64(0.0).Mul(&x.R, &y.R)
	yv := new(big.Float).SetFloat64(0.0).Mul(&x.I, &y.I)

	negXv := new(big.Float).SetFloat64(0.0).Mul(&x.R, &y.I)
	negXv.Neg(negXv)
	yu := new(big.Float).SetFloat64(0.0).Mul(&x.I, &y.R)

	dividendReal := new(big.Float).Add(xu, yv)
	dividendImag := new(big.Float).Add(negXv, yu)

	z.R = *new(big.Float).Quo(dividendReal, divisor)
	z.I = *new(big.Float).Quo(dividendImag, divisor)
	return z
}

// sqrt(x.Real^2 + x.Imag^2)
func Abs(x *BigComplex) *big.Float {
	realsquared := new(big.Float).Set(&x.R)
	realsquared.Mul(realsquared, realsquared)

	imagsquared := new(big.Float).Set(&x.I)
	imagsquared.Mul(imagsquared, imagsquared)

	realsquared.Add(realsquared, imagsquared)

	return realsquared.Sqrt(realsquared)
}

func (z *BigComplex) Set(x *BigComplex) {
	z.R.Set(&x.R)
	z.I.Set(&x.I)
}

func Real(b *BigComplex) *big.Float {
	return &b.R
}

func Imag(b *BigComplex) *big.Float {
	return &b.I
}

func SetComplex(c complex128) *BigComplex {
	z := new(BigComplex)
	z.R = *big.NewFloat(real(c))
	z.I = *big.NewFloat(imag(c))

	return z
}

func (z *BigComplex) Complex128() complex128 {
	real, _ := Real(z).Float64()
	imag, _ := Imag(z).Float64()

	return complex(real, imag)
}
