package main

import "fmt"
import "testing"

// Returns true if p and q are the same size and have equal entries.
func wordArraysEq(p, q []Word) bool {
	if len(p) != len(q) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

// Dumps p to a string.
func dumpWordPoly(p *WordPoly) string {
	s := ""
	for i := len(p.coeffs) - 1; i >= 0; i-- {
		if p.coeffs[i] > 0 {
			if s != "" {
				s += " + "
			}
			s += fmt.Sprintf("%dx^%d", p.coeffs[i], i)
		}
	}
	if s == "" {
		return "0"
	}
	return s
}

// NewWordPoly(k, a, N, R) should return a polynomial representing
// X^(k % R) + (a % N).
func TestNewWordPoly(t *testing.T) {
	var a Word = 12
	var k Word = 6
	var N Word = 10
	var R Word = 5
	p := NewWordPoly(a, k, N, R)
	q := []Word{2, 1, 0, 0, 0}
	if !wordArraysEq(p.coeffs, q) {
		t.Error(dumpWordPoly(p))
	}
}

// p.Eq(q) should return whether p and q have the same coefficients.
func TestWordPolyEq(t *testing.T) {
	var N Word = 10
	var R Word = 5

	p := NewWordPoly(1, 2, N, R)
	q := NewWordPoly(1, 3, N, R)
	r := NewWordPoly(2, 3, N, R)

	// Test reflexivity.
	if !p.Eq(p) {
		t.Error(dumpWordPoly(p))
	}
	if !q.Eq(q) {
		t.Error(dumpWordPoly(q))
	}
	if !r.Eq(r) {
		t.Error(dumpWordPoly(r))
	}

	if p.Eq(q) {
		t.Error(dumpWordPoly(p), dumpWordPoly(q))
	}
	if p.Eq(r) {
		t.Error(dumpWordPoly(p), dumpWordPoly(r))
	}
	if q.Eq(p) {
		t.Error(dumpWordPoly(q), dumpWordPoly(p))
	}
	if q.Eq(r) {
		t.Error(dumpWordPoly(q), dumpWordPoly(r))
	}
	if r.Eq(p) {
		t.Error(dumpWordPoly(r), dumpWordPoly(p))
	}
	if r.Eq(q) {
		t.Error(dumpWordPoly(r), dumpWordPoly(q))
	}
}

// Multiplication should be modulo (N, X^R - 1).
func TestWordPolyMul(t *testing.T) {
	var N Word = 10
	var R Word = 5

	p := NewWordPoly(4, 3, N, R)
	tmp := NewWordPoly(0, 0, N, R)
	p.mul(p, N, tmp)
	q := []Word{6, 1, 0, 8, 0}
	if !wordArraysEq(p.coeffs, q) {
		t.Error(dumpWordPoly(p))
	}
}

// (X + a)^N should equal X^n + a mod (N, X^R - 1) for prime N.
func TestWordPolyPow(t *testing.T) {
	var a Word = 2
	var N Word = 101
	var R Word = 53

	p := NewWordPoly(a, 1, N, R)
	tmp1 := NewWordPoly(0, 0, N, R)
	tmp2 := NewWordPoly(0, 0, N, R)
	p.Pow(N, tmp1, tmp2)
	q := NewWordPoly(a, N, N, R)
	if !wordArraysEq(p.coeffs, q.coeffs) {
		t.Error(dumpWordPoly(p), dumpWordPoly(q))
	}
}

// Make sure that polynomials get converted to strings in standard
// notation.
func TestWordPolyFormat(t *testing.T) {
	var N Word = 101
	var R Word = 53

	p := &WordPoly{}
	str := fmt.Sprint(p)
	if str != "0" {
		t.Error(dumpWordPoly(p), str)
	}

	p = NewWordPoly(2, 3, N, R)
	str = fmt.Sprint(p)
	if str != "x^3 + 2" {
		t.Error(dumpWordPoly(p), str)
	}

	p = NewWordPoly(1, 1, N, R)
	str = fmt.Sprint(p)
	if str != "x + 1" {
		t.Error(dumpWordPoly(p), str)
	}
}
