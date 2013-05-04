package main

import "fmt"
import "math/big"
import "testing"

// Fill p's unused bits with non-zero data. This helps in flushing out
// any bugs related to relying on memory to be zeroed.
func fuzzBigIntPoly(p *BigIntPoly) {
	bits := p.phi.Bits()
	unusedBits := bits[len(bits):cap(bits)]
	for i := 0; i < len(unusedBits); i++ {
		unusedBits[i] = ^big.Word(0)
	}
}

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
func bigIntPolyHasCoefficients(p *BigIntPoly, coefficients []int64) bool {
	e := calculatePhi(coefficients, p.k)
	return p.phi.Cmp(&e) == 0
}

// Dumps p to a string.
func dumpBigIntPoly(p *BigIntPoly) string {
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

// NewBigIntPoly(k, a, N, R) should return the zero polynomial
// mod (N, X^R - 1).
func TestNewBigIntPoly(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)
	p := NewBigIntPoly(N, R)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasCoefficients(p, []int64{}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// BigIntPoly.Set() should set the polynomial to X^(k % R) + (a % N).
func TestBigIntPolySet(t *testing.T) {
	a := *big.NewInt(12)
	k := *big.NewInt(6)
	N := *big.NewInt(10)
	R := *big.NewInt(5)
	p := NewBigIntPoly(N, R)
	p.Set(a, k, N)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasCoefficients(p, []int64{2, 1}) {
		t.Error(dumpBigIntPoly(p))
	}

	a = *big.NewInt(13)
	k = *big.NewInt(7)
	p.Set(a, k, N)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasCoefficients(p, []int64{3, 0, 1}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// p.Eq(q) should return whether p and q have the same coefficients.
func TestBigIntPolyEq(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)

	p := NewBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(2), N)
	fuzzBigIntPoly(p)
	q := NewBigIntPoly(N, R)
	q.Set(*big.NewInt(1), *big.NewInt(3), N)
	fuzzBigIntPoly(q)
	r := NewBigIntPoly(N, R)
	r.Set(*big.NewInt(2), *big.NewInt(3), N)
	fuzzBigIntPoly(r)

	// Test reflexivity.
	if !p.Eq(p) {
		t.Error(dumpBigIntPoly(p))
	}
	if !q.Eq(q) {
		t.Error(dumpBigIntPoly(q))
	}
	if !r.Eq(r) {
		t.Error(dumpBigIntPoly(r))
	}

	if p.Eq(q) {
		t.Error(dumpBigIntPoly(p), dumpBigIntPoly(q))
	}
	if p.Eq(r) {
		t.Error(dumpBigIntPoly(p), dumpBigIntPoly(r))
	}
	if q.Eq(p) {
		t.Error(dumpBigIntPoly(q), dumpBigIntPoly(p))
	}
	if q.Eq(r) {
		t.Error(dumpBigIntPoly(q), dumpBigIntPoly(r))
	}
	if r.Eq(p) {
		t.Error(dumpBigIntPoly(r), dumpBigIntPoly(p))
	}
	if r.Eq(q) {
		t.Error(dumpBigIntPoly(r), dumpBigIntPoly(q))
	}
}

// Multiplication should be modulo (N, X^R - 1).
func TestBigIntPolyMul(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)

	p := NewBigIntPoly(N, R)
	p.Set(*big.NewInt(4), *big.NewInt(3), N)
	fuzzBigIntPoly(p)
	tmp := NewBigIntPoly(N, R)
	fuzzBigIntPoly(tmp)
	p.mul(p, N, tmp)
	if !bigIntPolyHasCoefficients(p, []int64{6, 1, 0, 8}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// (X + a)^N should equal X^n + a mod (N, X^R - 1) for prime N.
func TestBigIntPolyPow(t *testing.T) {
	a := *big.NewInt(2)
	N := *big.NewInt(101)
	R := *big.NewInt(53)

	p := NewBigIntPoly(N, R)
	p.Set(a, *big.NewInt(1), N)
	fuzzBigIntPoly(p)
	tmp1 := NewBigIntPoly(N, R)
	tmp2 := NewBigIntPoly(N, R)
	fuzzBigIntPoly(tmp1)
	fuzzBigIntPoly(tmp2)
	p.Pow(N, tmp1, tmp2)
	q := NewBigIntPoly(N, R)
	q.Set(a, N, N)
	fuzzBigIntPoly(q)
	if p.phi.Cmp(&q.phi) != 0 {
		t.Error(dumpBigIntPoly(p), dumpBigIntPoly(q))
	}
}

// Make sure that polynomials get converted to strings in standard
// notation.
func TestBigIntPolyFormat(t *testing.T) {
	N := *big.NewInt(101)
	R := *big.NewInt(53)

	p := &BigIntPoly{}
	fuzzBigIntPoly(p)
	str := fmt.Sprint(p)
	if str != "0" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = NewBigIntPoly(N, R)
	p.Set(*big.NewInt(2), *big.NewInt(3), N)
	fuzzBigIntPoly(p)
	str = fmt.Sprint(p)
	if str != "x^3 + 2" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = NewBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(1), N)
	fuzzBigIntPoly(p)
	str = fmt.Sprint(p)
	if str != "x + 1" {
		t.Error(dumpBigIntPoly(p), str)
	}
}
