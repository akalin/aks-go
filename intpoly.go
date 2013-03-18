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

// An IntPoly represents the polynomial with the given non-zero terms
// in order of ascending degree.
// The zero value for an IntPoly represents the zero polynomial.
type IntPoly struct {
	terms []intMono
}

// Builds a new IntPoly from the given list of coefficient/degree
// pairs. Each coefficient must be non-zero, each degree must be
// non-negative, and the list must be in ascending order of degree.
func NewIntPoly(terms [][2]*big.Int) *IntPoly {
	p := IntPoly{make([]intMono, len(terms))}
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
