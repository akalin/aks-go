package main

import "fmt"
import "math/big"
import "testing"

// Returns whether p is the zero polynomial.
func isZero(p *IntPoly) bool {
	return len(p.terms) == 0
}

// Converts a list of int64 pairs to a list of *big.Int pairs.
func makeTerms(int64Terms [][2]int64) [][2]*big.Int {
	terms := make([][2]*big.Int, len(int64Terms))
	for i, int64Term := range int64Terms {
		terms[i][0] = big.NewInt(int64Term[0])
		terms[i][1] = big.NewInt(int64Term[1])
	}
	return terms
}

// Returns whether p has exactly the given terms.
func hasTerms(p *IntPoly, terms [][2]int64) bool {
	if len(p.terms) != len(terms) {
		return false
	}
	for i, term := range p.terms {
		if term.coeff.Cmp(big.NewInt(terms[i][0])) != 0 ||
			term.deg.Cmp(big.NewInt(terms[i][1])) != 0 {
			return false
		}
	}
	return true
}

// Dumps p to a string.
func dumpIntPoly(p *IntPoly) string {
	if isZero(p) {
		return "0"
	}
	s := ""
	for i := len(p.terms) - 1; i >= 0; i-- {
		s += fmt.Sprintf("%vx^%v", &p.terms[i].coeff, &p.terms[i].deg)
		if i > 0 {
			s += " + "
		}
	}
	return s
}

// Passing an empty slice to NewIntPoly() should give a zero polynomial.
func TestNewIntPolyEmpty(t *testing.T) {
	p := NewIntPoly([][2]*big.Int{})
	if !isZero(p) {
		t.Error(dumpIntPoly(p))
	}
}

// NewIntPoly() should return a polynomial with the given terms.
func TestNewIntPolyBasic(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {6, 6}, {-7, 9}}
	p := NewIntPoly(makeTerms(terms))
	if !hasTerms(p, terms) {
		t.Error(dumpIntPoly(p))
	}
}

// Eq() should return true iff its given polynomials have the same
// terms.
func TestIntPolyEq(t *testing.T) {
	terms1 := [][2]int64{{1, 1}, {-2, 4}, {3, 6}, {-7, 9}}
	terms2 := [][2]int64{{2, 1}, {-2, 4}, {3, 6}, {-7, 9}}
	p1 := NewIntPoly(makeTerms(terms1))
	p2 := NewIntPoly(makeTerms(terms2))
	p3 := NewIntPoly(makeTerms(terms2[0:3]))

	// Reflexivity.
	if !p1.Eq(p1) {
		t.Error(dumpIntPoly(p1))
	}
	if !p2.Eq(p2) {
		t.Error(dumpIntPoly(p2))
	}
	if !p3.Eq(p3) {
		t.Error(dumpIntPoly(p3))
	}

	// Symmetry.
	if p1.Eq(p2) != p2.Eq(p1) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p2))
	}
	if p1.Eq(p3) != p3.Eq(p1) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p3))
	}
	if p2.Eq(p3) != p3.Eq(p2) {
		t.Error(dumpIntPoly(p2), dumpIntPoly(p3))
	}

	// Transitivity.
	p4 := NewIntPoly(makeTerms(terms1))
	p5 := NewIntPoly(makeTerms(terms1))
	if !p1.Eq(p4) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p4))
	}
	if !p4.Eq(p5) {
		t.Error(dumpIntPoly(p4), dumpIntPoly(p5))
	}
	if !p1.Eq(p5) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p5))
	}

	// p1 and p2 don't have the same coefficient.
	if p1.Eq(p2) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p2))
	}

	// p1 and p3 don't have the same degree.
	if p1.Eq(p3) {
		t.Error(dumpIntPoly(p1), dumpIntPoly(p3))
	}
}
