package main

import "fmt"

// A Word represents a coefficient of a WordPoly.
// TODO(akalin): Use uintptr instead.
type Word uint32

// The size of Word in bits.
const WORD_BITS = 32

// A WordPoly represents a polynomial with Word coefficients.
//
// The zero value for a WordPoly represents the zero polynomial.
type WordPoly struct {
	coeffs []Word
}

// Only polynomials built with the same value of N and R may be used
// together in one of the functions below.

// Builds a new WordPoly representing the zero polynomial
// mod (N, X^R - 1). R must fit into an int.
func NewWordPoly(N, R Word) *WordPoly {
	return &WordPoly{make([]Word, R)}
}

// Sets p to X^k + a mod (N, X^R - 1).
func (p *WordPoly) Set(a, k, N Word) {
	R := len(p.coeffs)
	p.coeffs[0] = a % N
	for i := 1; i < R; i++ {
		p.coeffs[i] = 0
	}
	p.coeffs[int(k%Word(R))] = 1
}

// Returns whether p has the same coefficients as q.
func (p *WordPoly) Eq(q *WordPoly) bool {
	R := len(p.coeffs)
	for i := 0; i < R; i++ {
		if p.coeffs[i] != q.coeffs[i] {
			return false
		}
	}
	return true
}

// Sets p to the product of p and q mod (N, X^R - 1). tmp must not
// alias p or q.
func (p *WordPoly) mul(q *WordPoly, N Word, tmp *WordPoly) {
	R := len(tmp.coeffs)
	for i := 0; i < R; i++ {
		tmp.coeffs[i] = 0
	}

	// Optimized and unrolled version of the following loop:
	//
	//   for i, j < R {
	//     tmp_{(i + j) % R} += (p_i * q_j) % N
	//   }
	for i := 0; i < R; i++ {
		for j := 0; j < R-i; j++ {
			k := i + j
			// TODO(akalin): Handle overflow here when we
			// change Word to uintptr.
			e := uint64(p.coeffs[i]) * uint64(q.coeffs[j])
			e %= uint64(N)
			e += uint64(tmp.coeffs[k])
			e %= uint64(N)
			tmp.coeffs[k] = Word(e)
		}
		for j := R - i; j < R; j++ {
			k := j - (R - i)
			// Duplicate of loop above.
			e := uint64(p.coeffs[i]) * uint64(q.coeffs[j])
			e %= uint64(N)
			e += uint64(tmp.coeffs[k])
			e %= uint64(N)
			tmp.coeffs[k] = Word(e)
		}
	}

	p.coeffs, tmp.coeffs = tmp.coeffs, p.coeffs
}

// Sets p to p^N mod (N, X^R - 1), where R is the size of p. N must be
// positive, and tmp1 and tmp2 must not alias each other or p.
func (p *WordPoly) Pow(N Word, tmp1, tmp2 *WordPoly) {
	R := len(p.coeffs)
	for i := 0; i < R; i++ {
		tmp1.coeffs[i] = p.coeffs[i]
	}

	// Find N's highest set bit.
	i := WORD_BITS - 1
	for ; (i >= 0) && ((N & (1 << uint(i))) == 0); i-- {
	}

	for i--; i >= 0; i-- {
		tmp1.mul(tmp1, N, tmp2)
		if (N & (1 << uint(i))) != 0 {
			tmp1.mul(p, N, tmp2)
		}
	}
	p.coeffs, tmp1.coeffs = tmp1.coeffs, p.coeffs
}

// fmt.Formatter implementation.
func (p *WordPoly) Format(f fmt.State, c rune) {
	i := len(p.coeffs) - 1
	for ; i >= 0 && p.coeffs[i] == 0; i-- {
	}

	if i < 0 {
		fmt.Fprint(f, "0")
		return
	}

	// Formats coeff*x^deg.
	formatNonZeroMonomial := func(f fmt.State, c rune, coeff, deg Word) {
		if coeff != 1 || deg == 0 {
			fmt.Fprint(f, coeff)
		}
		if deg != 0 {
			fmt.Fprint(f, "x")
			if deg > 1 {
				fmt.Fprint(f, "^", deg)
			}
		}
	}

	formatNonZeroMonomial(f, c, p.coeffs[i], Word(i))

	for i--; i >= 0; i-- {
		if p.coeffs[i] != 0 {
			fmt.Fprint(f, " + ")
			formatNonZeroMonomial(f, c, p.coeffs[i], Word(i))
		}
	}
}
