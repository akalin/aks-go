package main

import "math/big"

const (
	// Compute the size of a big.Word in bits.
	_m             = ^big.Word(0)
	_logS          = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S             = 1 << _logS
	_BIG_WORD_BITS = _S << 3
)

// A BigIntPoly2 represents a polynomial with big.Int coefficients mod
// some (N, X^R - 1).
//
// The zero value for a BigIntPoly represents the zero polynomial.
type BigIntPoly2 struct {
	R int
	// k is the number of big.Words required to hold a coefficient
	// in calculations without overflowing.
	k int
	// If p(x) is the polynomial as a function, phi is
	// p(2^{k*_BIG_WORD_BITS}). Since a coefficient fits into k
	// big.Words, this is a lossless transformation; that is, one
	// can recover all coefficients of p(x) from phi.
	phi big.Int
}

// Only polynomials built with the same value of N and R may be used
// together in one of the functions below.

// Builds a new BigIntPoly representing the zero polynomial
// mod (N, X^R - 1). R must fit into an int.
func NewBigIntPoly2(N, R big.Int) *BigIntPoly2 {
	// A coefficient can be up to R*(N - 1)^2 in intermediate
	// calculations.
	var maxCoefficient big.Int
	maxCoefficient.Sub(&N, big.NewInt(1))
	maxCoefficient.Mul(&maxCoefficient, &maxCoefficient)
	maxCoefficient.Mul(&maxCoefficient, &R)

	var phi big.Int
	rInt := int(R.Int64())
	k := len(maxCoefficient.Bits())
	// Up to 2*R coefficients may be needed in intermediate
	// calculations.
	maxWordCount := 2 * rInt * k
	phi.SetBits(make([]big.Word, maxWordCount))
	return &BigIntPoly2{rInt, k, phi}
}

// Returns 1 + the degree of this polynomial, or 0 if the polynomial
// is the zero polynomial.
func (p *BigIntPoly2) getCoefficientCount() int {
	l := len(p.phi.Bits())
	if l == 0 {
		return 0
	}
	coefficientCount := l / p.k
	if l%p.k != 0 {
		coefficientCount++
	}
	return coefficientCount
}

// Returns the ith coefficient of this polynomial. i must be less than
// p.getCoefficientCount().
func (p *BigIntPoly2) getCoefficient(i int) big.Int {
	var c big.Int
	start := i * p.k
	if i == p.getCoefficientCount()-1 {
		// If the last coefficient is small enough, phi might
		// have fewer than p.R * p.k words.
		c.SetBits(p.phi.Bits()[start:])
	} else {
		end := (i + 1) * p.k
		c.SetBits(p.phi.Bits()[start:end])
	}
	return c
}
