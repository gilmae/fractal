package fractal

import "math/big"

type BigComplex struct {
	r, i big.Float
}

/* Multiplying two complex numbers looks like a quadratic equation
 *
 * (x + yi)(w + zi) == w * w + xzi + wyi + yzi^2
 * And since i^2 = -1, this becomes ww-yz + (xz+wy)i
 */
 func (z *BigComplex) Mul(x, y *BigComplex) *BigComplex {
	rSquared := new(big.Float).Copy(&x.r)

	rSquared.Mul(rSquared, &y.r)

	iSquared := new(big.Float).Copy(&x.i)
	iSquared.Mul(iSquared, &y.i)

	r := new(big.Float).Copy(rSquared)
	r.Sub(r, iSquared)

	i := new(big.Float).Copy(&x.r)
	i.Mul(i, &y.i)

	tmp := new(big.Float).Copy(&x.i)
	tmp.Mul(tmp, &y.r)

	i.Add(i, tmp)

	z.r = *r
	z.i = *i
	
	return z
}

func (z *BigComplex) Add(x,y *BigComplex) *BigComplex {
	z.r.Add(&x.r, &y.r)
	z.i.Add(&x.i, &y.i)
	
	return z
}

func (z *BigComplex) Sub(x,y *BigComplex) *BigComplex {
	z.r.Sub(&x.r, &y.r)
	z.i.Sub(&x.i, &y.i)
	
	return z
}

/* 
	x + yi / u + vi = (xu + yv) + (-xv + yu)i / u^2 + v^2
 */
 func (z *BigComplex) Quo(x,y *BigComplex) *BigComplex {
	divisor := new(big.Float).Set(&y.r)
	divisor.Mul(divisor, divisor) // u^2
	tmp := new(big.Float).Set(&y.i)
	tmp.Mul(tmp,tmp) // v^2
	divisor.Add(divisor, tmp) // u^2 + v^2
	
	
	
	xu := new(big.Float).SetFloat64(0.0).Mul(&x.r, &y.r)
	yv := new(big.Float).SetFloat64(0.0).Mul(&x.i, &y.i)
	
	negXv := new(big.Float).SetFloat64(0.0).Mul(&x.r, &y.i)
	negXv.Neg(negXv)
	yu := new(big.Float).SetFloat64(0.0).Mul(&x.i, &y.r)
		
	dividend_real := new(big.Float).Add(xu,yv)
	dividend_imag := new(big.Float).Add(negXv,yu)
	
	
	z.r = *new(big.Float).Quo(dividend_real, divisor)
	z.i = *new(big.Float).Quo(dividend_imag, divisor)
	return z
}

// sqrt(x.real^2 + x.Imag^2)
func Abs(x *BigComplex) *big.Float{
	realsquared := new(big.Float).Set(&x.r)
	realsquared .Mul(realsquared , realsquared )	
	
	imagsquared := new(big.Float).Set(&x.i)
	imagsquared.Mul(imagsquared, imagsquared)
	
	realsquared.Add(realsquared, imagsquared)
	
	return realsquared.Sqrt(realsquared)
}