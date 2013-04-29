package main

import "fmt"
import "math/big"
import "testing"

func makeBigIntArray(ints []int64) []big.Int {
	bigInts := make([]big.Int, len(ints))
	for i := 0; i < len(ints); i++ {
		bigInts[i] = *big.NewInt(ints[i])
	}
	return bigInts
}

// Returns true if p and q are the same size and have equal entries.
func bigIntArraysEq(p, q []big.Int) bool {
	if len(p) != len(q) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i].Cmp(&q[i]) != 0 {
			return false
		}
	}
	return true
}

// Dumps p to a string.
func dumpBigIntPoly(p *BigIntPoly) string {
	s := ""
	for i := len(p.coeffs) - 1; i >= 0; i-- {
		if p.coeffs[i].Sign() > 0 {
			if s != "" {
				s += " + "
			}
			s += fmt.Sprintf("%vx^%d", &p.coeffs[i], i)
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
	q := makeBigIntArray([]int64{0, 0, 0, 0, 0})
	if !bigIntArraysEq(p.coeffs, q) {
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
	q := makeBigIntArray([]int64{2, 1, 0, 0, 0})
	if !bigIntArraysEq(p.coeffs, q) {
		t.Error(dumpBigIntPoly(p))
	}

	a = *big.NewInt(13)
	k = *big.NewInt(7)
	p.Set(a, k, N)
	q = makeBigIntArray([]int64{3, 0, 1, 0, 0})
	if !bigIntArraysEq(p.coeffs, q) {
		t.Error(dumpBigIntPoly(p))
	}
}

// p.Eq(q) should return whether p and q have the same coefficients.
func TestBigIntPolyEq(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)

	p := NewBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(2), N)
	q := NewBigIntPoly(N, R)
	q.Set(*big.NewInt(1), *big.NewInt(3), N)
	r := NewBigIntPoly(N, R)
	r.Set(*big.NewInt(2), *big.NewInt(3), N)

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
	tmp := NewBigIntPoly(N, R)
	p.mul(p, N, tmp)
	q := makeBigIntArray([]int64{6, 1, 0, 8, 0})
	if !bigIntArraysEq(p.coeffs, q) {
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
	tmp1 := NewBigIntPoly(N, R)
	tmp2 := NewBigIntPoly(N, R)
	p.Pow(N, tmp1, tmp2)
	q := NewBigIntPoly(N, R)
	q.Set(a, N, N)
	if !bigIntArraysEq(p.coeffs, q.coeffs) {
		t.Error(dumpBigIntPoly(p), dumpBigIntPoly(q))
	}
}

// Make sure that polynomials get converted to strings in standard
// notation.
func TestBigIntPolyFormat(t *testing.T) {
	N := *big.NewInt(101)
	R := *big.NewInt(53)

	p := &BigIntPoly{}
	str := fmt.Sprint(p)
	if str != "0" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = NewBigIntPoly(N, R)
	p.Set(*big.NewInt(2), *big.NewInt(3), N)
	str = fmt.Sprint(p)
	if str != "x^3 + 2" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = NewBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(1), N)
	str = fmt.Sprint(p)
	if str != "x + 1" {
		t.Error(dumpBigIntPoly(p), str)
	}
}
