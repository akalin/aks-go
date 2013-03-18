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

// Set() should make a deep copy.
func TestIntPolySetDeepCopy(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {3, 6}, {-7, 9}}
	p1 := NewIntPoly(makeTerms(terms))
	p2 := IntPoly{}
	p2Alias := p2.Set(p1)
	if &p2 != p2Alias {
		t.Errorf("%p %p", p2, p2Alias)
	}
	// This shouldn't affect the values in p2.
	p1.terms[0].coeff.SetInt64(2)
	if !hasTerms(&p2, terms) {
		t.Error(dumpIntPoly(&p2))
	}
}

// Setting a polynomial to itself should have no effect.
func TestIntPolySetSelf(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {3, 6}, {-7, 9}}
	p := NewIntPoly(makeTerms(terms))
	p.Set(p)
	if !hasTerms(p, terms) {
		t.Error(dumpIntPoly(p))
	}
}

// Add() should add its given polynomials term by term.
func TestIntPolyAdd(t *testing.T) {
	terms1 := [][2]int64{{1, 1}, {-2, 4}, {3, 6}}
	terms2 := [][2]int64{{2, 5}, {3, 6}, {-3, 7}}
	termsSum := [][2]int64{{1, 1}, {-2, 4}, {2, 5}, {6, 6}, {-3, 7}}
	p1 := NewIntPoly(makeTerms(terms1))
	p2 := NewIntPoly(makeTerms(terms2))
	sum := IntPoly{}
	sumAlias := sum.Add(p1, p2)
	if &sum != sumAlias {
		t.Errorf("%p %p", sum, sumAlias)
	}
	if !hasTerms(&sum, termsSum) {
		t.Error(dumpIntPoly(&sum))
	}
}

// Add() should still work even with aliasing.
func TestIntPolyAddAlias(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {3, 6}}
	p := NewIntPoly(makeTerms(terms))
	p.Add(p, p)
	if !hasTerms(p, [][2]int64{{2, 1}, {-4, 4}, {6, 6}}) {
		t.Error(dumpIntPoly(p))
	}
}

// MulMono() should multiply its given polynomial by its given
// monomial.
func TestIntPolyMulMono(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {3, 6}}
	p := NewIntPoly(makeTerms(terms))
	coeff := big.NewInt(5)
	deg := big.NewInt(3)

	termsProd := [][2]int64{{5, 4}, {-10, 7}, {15, 9}}
	prod := IntPoly{}
	prodAlias := prod.MulMono(p, coeff, deg)
	if &prod != prodAlias {
		t.Errorf("%p %p", prod, prodAlias)
	}
	if !hasTerms(&prod, termsProd) {
		t.Error(dumpIntPoly(&prod))
	}
}

// MulMono() should still work even with aliasing.
func TestIntPolyMulMonoAlias(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 4}, {3, 6}}
	p := NewIntPoly(makeTerms(terms))
	coeff := big.NewInt(5)
	deg := big.NewInt(3)

	termsProd := [][2]int64{{5, 4}, {-10, 7}, {15, 9}}
	p.MulMono(p, coeff, deg)
	if !hasTerms(p, termsProd) {
		t.Error(dumpIntPoly(p))
	}
}

// Mul() should multiply its given polynomials.
func TestIntPolyMul(t *testing.T) {
	terms1 := [][2]int64{{1, 1}, {-2, 3}, {3, 5}}
	terms2 := [][2]int64{{-2, 2}, {3, 4}}
	p := NewIntPoly(makeTerms(terms1))
	q := NewIntPoly(makeTerms(terms2))

	termsProd := [][2]int64{{-2, 3}, {7, 5}, {-12, 7}, {9, 9}}
	prod := IntPoly{}
	prodAlias := prod.Mul(p, q)
	if &prod != prodAlias {
		t.Errorf("%p %p", prod, prodAlias)
	}
	if !hasTerms(&prod, termsProd) {
		t.Error(dumpIntPoly(&prod))
	}
}

// Multiplication by the zero polynomial should result in the zero
// polynomial.
func TestIntPolyMulZero(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 3}}
	p := NewIntPoly(makeTerms(terms))
	prod := IntPoly{}
	prod.Mul(p, &IntPoly{})
	if !isZero(&prod) {
		t.Error(dumpIntPoly(&prod))
	}
}

// Mul() should still work even with aliasing.
func TestIntPolyMulAlias(t *testing.T) {
	terms := [][2]int64{{1, 1}, {-2, 3}}
	p := NewIntPoly(makeTerms(terms))

	termsProd := [][2]int64{{1, 2}, {-4, 4}, {4, 6}}
	p.Mul(p, p)

	if !hasTerms(p, termsProd) {
		t.Error(dumpIntPoly(p))
	}
}

// Pow() should raise its given polynomial by its given power.
func TestIntPolyPow(t *testing.T) {
	terms := [][2]int64{{1, 0}, {1, 1}}
	p := NewIntPoly(makeTerms(terms))

	termsPow := [][2]int64{{1, 0}, {4, 1}, {6, 2}, {4, 3}, {1, 4}}
	pow := IntPoly{}
	powAlias := pow.Pow(p, big.NewInt(4))
	if &pow != powAlias {
		t.Errorf("%p %p", pow, powAlias)
	}
	if !hasTerms(&pow, termsPow) {
		t.Error(dumpIntPoly(&pow))
	}
}

// Pow() should still work even with aliasing.
func TestIntPolyPowAlias(t *testing.T) {
	terms := [][2]int64{{1, 0}, {1, 1}}
	p := NewIntPoly(makeTerms(terms))

	termsPow := [][2]int64{{1, 0}, {3, 1}, {3, 2}, {1, 3}}
	p.Pow(p, big.NewInt(3))
	if !hasTerms(p, termsPow) {
		t.Error(dumpIntPoly(p))
	}
}

// Raising a non-zero polynomial to the zeroth power should give the
// constant polynomial 1.
func TestIntPolyNonZeroPowZero(t *testing.T) {
	terms := [][2]int64{{1, 0}, {1, 1}}
	p := NewIntPoly(makeTerms(terms))

	pow := IntPoly{}
	pow.Pow(p, big.NewInt(0))
	if !hasTerms(&pow, [][2]int64{{1, 0}}) {
		t.Error(dumpIntPoly(&pow))
	}
}

// Raising the zero polynomial to the zeroth power should give the
// the constant polynomial 1.
func TestIntPolyZeroPowZero(t *testing.T) {
	pow := IntPoly{}
	pow.Pow(&pow, big.NewInt(0))
	if !hasTerms(&pow, [][2]int64{{1, 0}}) {
		t.Error(dumpIntPoly(&pow))
	}
}
