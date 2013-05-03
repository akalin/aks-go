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

// Sets p to X^k + a mod (N, X^R - 1).
func (p *BigIntPoly2) Set(a, k, N big.Int) {
	R := big.NewInt(int64(p.R))
	var kModR big.Int
	kModR.Mod(&k, R)
	one := big.NewInt(1)
	p.phi.Lsh(one, uint(kModR.Int64())*uint(p.k*_BIG_WORD_BITS))

	var aModN big.Int
	aModN.Mod(&a, &N)
	p.phi.Add(&p.phi, &aModN)
}

// Returns whether p has the same coefficients as q.
func (p *BigIntPoly2) Eq(q *BigIntPoly2) bool {
	return p.phi.Cmp(&q.phi) == 0
}

// Sets p to the product of p and q mod (N, X^R - 1). tmp must not
// alias p or q.
func (p *BigIntPoly2) mul(q *BigIntPoly2, N big.Int, tmp *BigIntPoly2) {
	tmp.phi.Mul(&p.phi, &q.phi)

	// Mod tmp by X^R - 1.
	mid := p.R * p.k
	tmpBits := tmp.phi.Bits()
	if len(tmpBits) > mid {
		var lo, hi big.Int
		lo.SetBits(tmpBits[:mid])
		hi.SetBits(tmpBits[mid:])
		tmp.phi.Add(&lo, &hi)
	}

	// Set p to tmp mod N.
	p.phi.Set(&big.Int{})
	for i := tmp.getCoefficientCount() - 1; i >= 0; i-- {
		p.phi.Lsh(&p.phi, uint(p.k*_BIG_WORD_BITS))
		c := tmp.getCoefficient(i)
		if c.Cmp(&N) < 0 {
			p.phi.Add(&p.phi, &c)
		} else {
			// Mod c by N. Use big.Int.QuoRem() instead of
			// big.Int.Mod() since the latter allocates an
			// extra big.Int.
			c.QuoRem(&c, &N, &c)
			p.phi.Add(&p.phi, &c)
		}
	}
}
