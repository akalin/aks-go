package main

import "math/big"

// An intMono represents the monomial coeff*X^deg.
// The zero value for an intMono represents the zero monomial.
type intMono struct {
	coeff big.Int
	deg   big.Int
}

// Returns whether m and n have the same coefficient and degree.
func (m *intMono) Eq(n *intMono) bool {
	return m.coeff.Cmp(&n.coeff) == 0 && m.deg.Cmp(&n.deg) == 0
}

// Sets m to a deep copy of n.
func (m *intMono) Set(n *intMono) {
	m.coeff.Set(&n.coeff)
	m.deg.Set(&n.deg)
}

// The container for a polynomial's terms.
type termList []intMono

// Returns a termList of the given size, possibly reusing term's
// space.
func (terms termList) make(n int) termList {
	if n <= cap(terms) {
		// Reuse the space.
		return terms[0:n]
	}
	// Otherwise, allocate a new array. No need to worry about
	// allocating extra space since we don't have to worry about
	// carries.
	return make(termList, n)
}

// An IntPoly represents the polynomial with the given non-zero terms
// in order of ascending degree.
// The zero value for an IntPoly represents the zero polynomial.
type IntPoly struct {
	terms termList
}

// Builds a new IntPoly from the given list of coefficient/degree
// pairs. Each coefficient must be non-zero, each degree must be
// non-negative, and the list must be in ascending order of degree.
func NewIntPoly(terms [][2]*big.Int) *IntPoly {
	var p IntPoly
	p.terms = p.terms.make(len(terms))
	for i, term := range terms {
		if term[0].Sign() == 0 {
			panic("zero coefficient")
		}
		if term[1].Sign() < 0 {
			panic("negative degree")
		}
		if i > 0 && term[1].Cmp(terms[i-1][1]) <= 0 {
			panic("non-increasing degree")
		}
		p.terms[i].coeff.Set(term[0])
		p.terms[i].deg.Set(term[1])
	}
	return &p
}

// Returns whether p and q have the same terms.
func (p *IntPoly) Eq(q *IntPoly) bool {
	if len(p.terms) != len(q.terms) {
		return false
	}
	for i, pTerm := range p.terms {
		if !pTerm.Eq(&q.terms[i]) {
			return false
		}
	}
	return true
}

// Sets p to a deep copy of q.
func (p *IntPoly) Set(q *IntPoly) *IntPoly {
	if p == q {
		return p
	}
	p.terms = p.terms.make(len(q.terms))
	for i, _ := range p.terms {
		p.terms[i].Set(&q.terms[i])
	}
	return p
}
