package main

import "fmt"
import "math/big"
import "testing"

// Given a list of coefficients of a polynomial p(x) and the number of
// big.Words required to hold a coefficient, calculates phi =
// p(2^{k*_BIG_WORD_BITS}).
func calculatePhi(coefficients []int64, k int) big.Int {
	var e big.Int
	for i := len(coefficients) - 1; i >= 0; i-- {
		e.Lsh(&e, uint(k*_BIG_WORD_BITS))
		e.Add(&e, big.NewInt(coefficients[i]))
	}
	return e
}

// Returns whether p has exactly the given list of coefficients.
func bigIntPoly2HasCoefficients(p *BigIntPoly2, coefficients []int64) bool {
	e := calculatePhi(coefficients, p.k)
	return p.phi.Cmp(&e) == 0
}

// Dumps p to a string.
func dumpBigIntPoly2(p *BigIntPoly2) string {
	s := ""
	for i := p.getCoefficientCount() - 1; i >= 0; i-- {
		c := p.getCoefficient(i)
		if c.Sign() > 0 {
			if s != "" {
				s += " + "
			}
			s += fmt.Sprintf("%vx^%d", &c, i)
		}
	}
	if s == "" {
		return "0"
	}
	return s
}

// NewBigIntPoly2(k, a, N, R) should return the zero polynomial
// mod (N, X^R - 1).
func TestNewBigIntPoly2(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)
	p := NewBigIntPoly2(N, R)
	if !bigIntPoly2HasCoefficients(p, []int64{}) {
		t.Error(dumpBigIntPoly2(p))
	}
}
