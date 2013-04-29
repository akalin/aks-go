package main

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
	return &BigIntPoly{make([]big.Int, int(R.Int64()))}
}

// Sets p to X^k + a mod (N, X^R - 1).
func (p *BigIntPoly) Set(a, k, N big.Int) {
	R := len(p.coeffs)
	p.coeffs[0].Mod(&a, &N)
	for i := 1; i < R; i++ {
		p.coeffs[i] = big.Int{}
	}
	var i big.Int
	i.Mod(&k, big.NewInt(int64(R)))
	p.coeffs[int(i.Int64())] = *big.NewInt(1)
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
