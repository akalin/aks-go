package main

// A Word represents a coefficient of a WordPoly.
// TODO(akalin): Use uintptr instead.
type Word uint32

const _WORD_BITS = 32

// A WordPoly represents a polynomial with Word coefficients.
//
// The zero value for a WordPoly represents the zero polynomial.
type WordPoly struct {
	coeffs []Word
}

// Builds a new WordPoly representing X^k + a mod (N, X^R - 1). R must
// fit into an int.
func NewWordPoly(a, k, N, R Word) *WordPoly {
	p := WordPoly{make([]Word, R)}
	p.coeffs[0] = a % N
	p.coeffs[k%R] = 1
	return &p
}

// In the functions below, the polynomials in question must have been
// built with the same value of N and R.

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

	for i := 0; i < R; i++ {
		for j := 0; j < R; j++ {
			k := (i + j) % R
			// TODO(akalin): Handle overflow here when we
			// change Word to uintptr.
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
	i := _WORD_BITS - 1
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
