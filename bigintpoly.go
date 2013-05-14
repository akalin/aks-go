package main

import "fmt"
import "math/big"

// A BigIntPoly represents a polynomial with big.Int coefficients mod
// some (N, X^R - 1).
//
// The zero value for a BigIntPoly represents the zero polynomial.
type BigIntPoly struct {
	R int
	// k is the number of big.Words required to hold a coefficient
	// in calculations without overflowing.
	k int
	// If p(x) is the polynomial as a function, phi is
	// p(2^{k*bitsize(big.Word)}). Since a coefficient fits into k
	// big.Words, this is a lossless transformation; that is, one
	// can recover all coefficients of p(x) from phi.
	//
	// phi is set to have capacity for the largest possible
	// (intermediate) polynomial. No assumptions can be made about
	// the bytes in the unused capacity except for that the unused
	// bytes for the leading coefficient (if any) is guaranteed to
	// be zeroed out.
	phi big.Int
}

// Only polynomials built with the same value of N and R may be used
// together in one of the functions below.

// Builds a new BigIntPoly representing the zero polynomial
// mod (N, X^R - 1). R must fit into an int.
func NewBigIntPoly(N, R big.Int) *BigIntPoly {
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
	return &BigIntPoly{rInt, k, phi}
}

// Returns 1 + the degree of this polynomial, or 0 if the polynomial
// is the zero polynomial.
func (p *BigIntPoly) getCoefficientCount() int {
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

// Sets the coefficient count to the given number, which must be at
// most p.R. The unused bytes of the leading coefficient must be
// cleared (via commitCoefficient()) prior to this being called.
func (p *BigIntPoly) setCoefficientCount(coefficientCount int) {
	p.phi.SetBits(p.phi.Bits()[0 : coefficientCount*p.k])
}

// Returns the ith coefficient of this polynomial. i must be less than
// p.getCoefficientCount().
func (p *BigIntPoly) getCoefficient(i int) big.Int {
	var c big.Int
	start := i * p.k
	end := (i + 1) * p.k
	// Since the unused data for the leading coefficient is
	// guaranteed to be zeroed out, this is okay.
	c.SetBits(p.phi.Bits()[start:end])
	return c
}

// Clears the unused bytes of the given coefficient. Must be called
// after all changes have been made to a coefficient via a big.Int
// returned from p.getCoefficient(). Also must be called on the
// leading coefficient before p.setCoefficientCount() is called.
func (p *BigIntPoly) commitCoefficient(c big.Int) {
	cBits := c.Bits()
	unusedBits := cBits[len(cBits):p.k]
	for j := 0; j < len(unusedBits); j++ {
		unusedBits[j] = 0
	}
}

// Sets p to X^k + a mod (N, X^R - 1).
func (p *BigIntPoly) Set(a, k, N big.Int) {
	c0 := p.getCoefficient(0)
	c0.Mod(&a, &N)
	p.commitCoefficient(c0)

	R := big.NewInt(int64(p.R))
	var kModRBig big.Int
	kModRBig.Mod(&k, R)
	kModR := int(kModRBig.Int64())

	for i := 1; i <= kModR; i++ {
		c := p.getCoefficient(i)
		c.Set(&big.Int{})
		p.commitCoefficient(c)
	}

	cKModR := p.getCoefficient(kModR)
	cKModR.Set(big.NewInt(1))
	p.commitCoefficient(cKModR)

	p.setCoefficientCount(kModR + 1)
}

// Returns whether p has the same coefficients as q.
func (p *BigIntPoly) Eq(q *BigIntPoly) bool {
	return p.phi.Cmp(&q.phi) == 0
}

// Sets p to the product of p and q mod (N, X^R - 1). Assumes R >=
// 2. tmp must not alias p or q.
func (p *BigIntPoly) mul(
	q *BigIntPoly, N big.Int, tmp *BigIntPoly, prefix string) {
	fmt.Printf("%s: multiplying p and q\n", prefix)
	tmp.phi.Mul(&p.phi, &q.phi)
	fmt.Printf("%s: multiplying p and q done\n", prefix)
	p.phi, tmp.phi = tmp.phi, p.phi

	// Mod p by X^R - 1.
	mid := p.R * p.k
	pBits := p.phi.Bits()
	if len(pBits) > mid {
		fmt.Printf("%s: truncating\n", prefix)
		p.phi.SetBits(pBits[:mid])
		fmt.Printf("%s: truncating done\n", prefix)
	} else {
		fmt.Printf("%s: fits in lower half, not truncating\n", prefix)
	}
	/*
		// Clear the unused bits of the leading coefficient if
		// necessary.
		if len(pBits)%p.k != 0 {
			start := len(pBits)
			end := start + p.k - start%p.k
			unusedBits := pBits[start:end]
			for i := 0; i < len(unusedBits); i++ {
				unusedBits[i] = 0
			}
		}
		// Commit the leading coefficient before we access it.
		oldCoefficientCount := p.getCoefficientCount()
		if oldCoefficientCount > 0 {
			p.commitCoefficient(p.getCoefficient(oldCoefficientCount - 1))
		}
		fmt.Printf("%s: modding each coefficient by N\n", prefix)

		// Mod p by N.
		newCoefficientCount := 0
		tmp2 := tmp.getCoefficient(0)
		tmp3 := tmp.getCoefficient(1)
		for i := 0; i < oldCoefficientCount; i++ {
			c := p.getCoefficient(i)
			if c.Cmp(&N) >= 0 {
				// Mod c by N. Use big.Int.QuoRem() instead of
				// big.Int.Mod() since the latter allocates an
				// extra big.Int.
				tmp2.QuoRem(&c, &N, &tmp3)
				c.Set(&tmp3)
				p.commitCoefficient(c)
			}
			if c.Sign() != 0 {
				newCoefficientCount = i + 1
			}
		}
		p.setCoefficientCount(newCoefficientCount)

		fmt.Printf("%s: modding each coefficient by N done\n", prefix)
	*/
}

// Sets p to p^N mod (N, X^R - 1), where R is the size of p. tmp1 and
// tmp2 must not alias each other or p.
func (p *BigIntPoly) Pow(N big.Int, tmp1, tmp2 *BigIntPoly) {
	s := uint(p.R * p.k * _BIG_WORD_BITS)
	for i := 0; ; i++ {
		fmt.Printf("%d: multiplying...\n", i)
		p.phi.Mul(&p.phi, &p.phi)
		fmt.Printf("%d: multiplying done; shifting...\n", i)
		len := uint(p.phi.BitLen())
		if len > s {
			fmt.Printf("%d: shifting...\n", i)
			p.phi.Rsh(&p.phi, len - s)
			fmt.Printf("%d: shifting done.\n", i)
		} else {
			fmt.Printf("%d: not shifting\n", i)
		}
	}
}

// fmt.Formatter implementation.
func (p *BigIntPoly) Format(f fmt.State, c rune) {
	if p.phi.Sign() == 0 {
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

	i := p.getCoefficientCount() - 1
	formatNonZeroMonomial(f, c, p.getCoefficient(i), i)

	for i--; i >= 0; i-- {
		coeff := p.getCoefficient(i)
		if coeff.Sign() != 0 {
			fmt.Fprint(f, " + ")
			formatNonZeroMonomial(f, c, coeff, i)
		}
	}
}
