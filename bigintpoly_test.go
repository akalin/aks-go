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
