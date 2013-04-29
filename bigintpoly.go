package main

import "fmt"
import "math/big"

// A BigIntPoly represents a polynomial with big.Int coefficients.
//
// The zero value for a BigIntPoly represents the zero polynomial.
//
// TODO(akalin): Replace IntPoly with BigIntPoly.
type BigIntPoly struct {
	coeffs []big.Int
}

// Only polynomials built with the same value of N and R may be used
// together in one of the functions below.

// Builds a new BigIntPoly representing the zero polynomial
// mod (N, X^R - 1). R must fit into an int.
func NewBigIntPoly(N, R big.Int) *BigIntPoly {
	rInt := int(R.Int64())
	p := BigIntPoly{make([]big.Int, rInt)}

	// Pre-allocate space for each coefficient. In order to take
	// advantage of this, we must not assign to entries of
	// p.coeffs directly but instead use big.Int.Set().
	var nSq big.Int
	nSq.Mul(&N, &N)
	for i := 0; i < rInt; i++ {
		p.coeffs[i].Set(&nSq)
		p.coeffs[i].Set(&big.Int{})
	}

	return &p
}

// Sets p to X^k + a mod (N, X^R - 1).
func (p *BigIntPoly) Set(a, k, N big.Int) {
	R := len(p.coeffs)
	p.coeffs[0].Mod(&a, &N)
	for i := 1; i < R; i++ {
		p.coeffs[i].Set(&big.Int{})
	}
	var i big.Int
	i.Mod(&k, big.NewInt(int64(R)))
	p.coeffs[int(i.Int64())].Set(big.NewInt(1))
}

// Returns whether p has the same coefficients as q.
func (p *BigIntPoly) Eq(q *BigIntPoly) bool {
	R := len(p.coeffs)
	for i := 0; i < R; i++ {
		if p.coeffs[i].Cmp(&q.coeffs[i]) != 0 {
			return false
		}
	}
	return true
}

// Sets p to the product of p and q mod (N, X^R - 1). tmp1 and tmp2
// must not alias each other or p or q.
func (p *BigIntPoly) mul(q *BigIntPoly, N big.Int, tmp1, tmp2 *BigIntPoly) {
	R := len(tmp1.coeffs)
	for i := 0; i < R; i++ {
		tmp1.coeffs[i].Set(&big.Int{})
	}

	for i := 0; i < R; i++ {
		for j := 0; j < R; j++ {
			k := (i + j) % R
			// Set tmp1.coeffs[k] to (tmp1.coeffs[k] +
			// p.coeffs[i] * q.coeffs[j]) % N.

			tmp2.coeffs[k].Mul(&p.coeffs[i], &q.coeffs[j])

			// Set tmp2.coeffs[k] to tmp2.coeffs[k] +
			// tmp1.coeffs[k], avoid copying if possible.
			if tmp2.coeffs[k].Sign() == 0 {
				tmp2.coeffs[k], tmp1.coeffs[k] =
					tmp1.coeffs[k], tmp2.coeffs[k]
			} else if tmp1.coeffs[k].Sign() != 0 {
				tmp2.coeffs[k].Add(
					&tmp2.coeffs[k], &tmp1.coeffs[k])
			}

			// Use big.Int.QuoRem() instead of
			// big.Int.Mod() since the latter allocates an
			// extra big.Int.
			tmp2.coeffs[k].QuoRem(
				&tmp2.coeffs[k], &N, &tmp1.coeffs[k])
		}
	}
	p.coeffs, tmp1.coeffs = tmp1.coeffs, p.coeffs
}

// Sets p to p^N mod (N, X^R - 1), where R is the size of p. tmp1,
// tmp2, and tmp3 must not alias each other or p.
func (p *BigIntPoly) Pow(N big.Int, tmp1, tmp2, tmp3 *BigIntPoly) {
	R := len(p.coeffs)
	for i := 0; i < R; i++ {
		tmp1.coeffs[i].Set(&p.coeffs[i])
	}

	for i := N.BitLen() - 2; i >= 0; i-- {
		tmp1.mul(tmp1, N, tmp2, tmp3)
		if N.Bit(i) != 0 {
			tmp1.mul(p, N, tmp2, tmp3)
		}
	}
	p.coeffs, tmp1.coeffs = tmp1.coeffs, p.coeffs
}

// fmt.Formatter implementation.
func (p *BigIntPoly) Format(f fmt.State, c rune) {
	i := len(p.coeffs) - 1
	for ; i >= 0 && p.coeffs[i].Sign() == 0; i-- {
	}

	if i < 0 {
		fmt.Fprint(f, "0")
		return
	}

	// Formats coeff*x^deg.
	formatNonZeroMonomial := func(
		f fmt.State, c rune,
		coeff big.Int, deg int) {
		if coeff.Cmp(big.NewInt(1)) != 0 || deg == 0 {
			fmt.Fprint(f, &coeff)
		}
		if deg != 0 {
			fmt.Fprint(f, "x")
			if deg > 1 {
				fmt.Fprint(f, "^", deg)
			}
		}
	}

	formatNonZeroMonomial(f, c, p.coeffs[i], i)

	for i--; i >= 0; i-- {
		if p.coeffs[i].Sign() != 0 {
			fmt.Fprint(f, " + ")
			formatNonZeroMonomial(f, c, p.coeffs[i], i)
		}
	}
}
