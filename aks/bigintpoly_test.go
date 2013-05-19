package aks

import "fmt"
import "math/big"
import "testing"

const (
	// Compute the size of a big.Word in bits.
	_m             = ^big.Word(0)
	_logS          = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S             = 1 << _logS
	_BIG_WORD_BITS = _S << 3
)

// Fill p's unused bits with non-zero data. This helps in flushing out
// any bugs related to relying on memory to be zeroed.
func fuzzBigIntPoly(p *bigIntPoly) {
	bits := p.phi.Bits()
	unusedBits := bits[len(bits):cap(bits)]
	for i := 0; i < len(unusedBits); i++ {
		unusedBits[i] = ^big.Word(0)
	}
}

// Given a list of coefficients of a polynomial p(x) and the number of
// big.Words required to hold a coefficient, calculates phi =
// p(2^{k*_BIG_WORD_BITS}).
func calculatePhi(coefficients []big.Int, k int) big.Int {
	var e big.Int
	for i := len(coefficients) - 1; i >= 0; i-- {
		e.Lsh(&e, uint(k*_BIG_WORD_BITS))
		e.Add(&e, &coefficients[i])
	}
	return e
}

// Returns whether p has exactly the given list of coefficients.
func bigIntPolyHasCoefficients(p *bigIntPoly, coefficients []big.Int) bool {
	e := calculatePhi(coefficients, p.k)
	return p.phi.Cmp(&e) == 0
}

// Returns whether p has exactly the given list of int64 coefficients.
func bigIntPolyHasInt64Coefficients(
	p *bigIntPoly, int64Coefficients []int64) bool {
	coefficients := make([]big.Int, len(int64Coefficients))
	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = *big.NewInt(int64Coefficients[i])
	}
	return bigIntPolyHasCoefficients(p, coefficients)
}

// Dumps p to a string.
func dumpBigIntPoly(p *bigIntPoly) string {
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

// newBigIntPoly(k, a, N, R) should return the zero polynomial
// mod (N, X^R - 1).
func TestNewBigIntPoly(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)
	p := newBigIntPoly(N, R)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasInt64Coefficients(p, []int64{}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// bigIntPoly.Set() should set the polynomial to X^(k % R) + (a % N).
func TestBigIntPolySet(t *testing.T) {
	a := *big.NewInt(12)
	k := *big.NewInt(6)
	N := *big.NewInt(10)
	R := *big.NewInt(5)
	p := newBigIntPoly(N, R)
	p.Set(a, k, N)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasInt64Coefficients(p, []int64{2, 1}) {
		t.Error(dumpBigIntPoly(p))
	}

	a = *big.NewInt(13)
	k = *big.NewInt(7)
	p.Set(a, k, N)
	fuzzBigIntPoly(p)
	if !bigIntPolyHasInt64Coefficients(p, []int64{3, 0, 1}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// p.Eq(q) should return whether p and q have the same coefficients.
func TestBigIntPolyEq(t *testing.T) {
	N := *big.NewInt(10)
	R := *big.NewInt(5)

	p := newBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(2), N)
	fuzzBigIntPoly(p)
	q := newBigIntPoly(N, R)
	q.Set(*big.NewInt(1), *big.NewInt(3), N)
	fuzzBigIntPoly(q)
	r := newBigIntPoly(N, R)
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

	// p = X^3 + 4.
	p := newBigIntPoly(N, R)
	p.Set(*big.NewInt(4), *big.NewInt(3), N)
	fuzzBigIntPoly(p)
	tmp := newBigIntPoly(N, R)
	fuzzBigIntPoly(tmp)
	// p^2 = (X^3 + 4)^2 = X^6 + 8X^3 + 16 which should be equal
	// to 8X^3 + X + 6 mod (10, X^5 - 1).
	p.mul(p, N, tmp)
	if !bigIntPolyHasInt64Coefficients(p, []int64{6, 1, 0, 8}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// Multiplication should still work for large (multi-word) values of
// N.
func TestBigIntPolyMulLarge(t *testing.T) {
	// Set word size to 5.
	one := big.NewInt(1)
	var N big.Int
	N.Lsh(one, 2*_BIG_WORD_BITS)
	var R big.Int
	R.Lsh(one, 10)
	rInt := int(R.Int64())

	// p = X^{N-1} + (N-1).
	p := newBigIntPoly(N, R)
	if p.k != 5 {
		t.Error(dumpBigIntPoly(p))
	}
	var nMinusOne big.Int
	nMinusOne.Sub(&N, one)
	p.Set(nMinusOne, nMinusOne, N)
	fuzzBigIntPoly(p)

	// p^2 = (X^{N-1} + (N-1))^2 = X^{2(N-1)} + 2(N-1) + (N-1)^2,
	// which should be equal to (N-2)X^{R-1} + X^{R-2} + 1. (The
	// div/mod operations should put their results in-place.)
	tmp := newBigIntPoly(N, R)
	fuzzBigIntPoly(tmp)
	p.mul(p, N, tmp)

	coeffs := make([]big.Int, rInt)
	coeffs[0].Set(one)
	coeffs[rInt-2].Set(one)
	coeffs[rInt-1].Sub(&N, big.NewInt(2))
	if !bigIntPolyHasCoefficients(p, coeffs) {
		t.Error(dumpBigIntPoly(p))
	}
}

// Multiplication should handle the leading coefficient correctly.
func TestBigIntPolyMulLeadingCoefficient(t *testing.T) {
	// Set word size to 2.
	one := big.NewInt(1)
	var nSize uint = _BIG_WORD_BITS - 6
	var N big.Int
	N.Lsh(one, nSize)
	var R big.Int
	R.Lsh(one, 10)

	// p = sqrt(N)X^{R/2}.
	p := newBigIntPoly(N, R)
	if p.k != 2 {
		t.Error(p.k)
	}
	fuzzBigIntPoly(p)
	var sqrtN big.Int
	sqrtN.Lsh(one, nSize/2)
	var rHalf big.Int
	rHalf.Rsh(&R, 1)
	p.Set(big.Int{}, rHalf, N)
	leadingCoeff := p.getCoefficient(int(rHalf.Int64()))
	leadingCoeff.Set(&sqrtN)
	fuzzBigIntPoly(p)

	// p^2 = NX^R, which should be equal to 0 mod (N, R).
	tmp := newBigIntPoly(N, R)
	fuzzBigIntPoly(tmp)
	p.mul(p, N, tmp)

	if !bigIntPolyHasCoefficients(p, []big.Int{}) {
		t.Error(dumpBigIntPoly(p))
	}
}

// Multiplication should handle the unused bytes of the leading
// coefficient correctly.
func TestBigIntPolyMulLeadingCoefficientUnusedBytes(t *testing.T) {
	// Set word size to 3.
	one := big.NewInt(1)
	var N big.Int
	N.Lsh(one, _BIG_WORD_BITS)
	// Set N to something not too close to a power of 2 to avoid
	// masking bugs.
	N.Sub(&N, big.NewInt(5))
	var R big.Int
	R.Lsh(one, 10)
	R.Sub(&R, one)

	// p = X, which should take up 4 words.
	p := newBigIntPoly(N, R)
	if p.k != 3 {
		t.Error(dumpBigIntPoly(p))
	}
	fuzzBigIntPoly(p)
	p.Set(big.Int{}, *one, N)
	fuzzBigIntPoly(p)

	// p^2 = X^2, which should take up 7 words. The unused 2 words
	// for the leading coefficient should not affect the result of
	// the multiplication.
	tmp := newBigIntPoly(N, R)
	fuzzBigIntPoly(tmp)
	p.mul(p, N, tmp)

	coeffs := []big.Int{big.Int{}, big.Int{}, *one}
	if !bigIntPolyHasCoefficients(p, coeffs) {
		t.Error(dumpBigIntPoly(p))
	}
}

// (X + a)^N should equal X^n + a mod (N, X^R - 1) for prime N.
func TestBigIntPolyPow(t *testing.T) {
	a := *big.NewInt(2)
	N := *big.NewInt(101)
	R := *big.NewInt(53)

	p := newBigIntPoly(N, R)
	p.Set(a, *big.NewInt(1), N)
	fuzzBigIntPoly(p)
	tmp1 := newBigIntPoly(N, R)
	tmp2 := newBigIntPoly(N, R)
	fuzzBigIntPoly(tmp1)
	fuzzBigIntPoly(tmp2)
	p.Pow(N, tmp1, tmp2)
	q := newBigIntPoly(N, R)
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

	p := &bigIntPoly{}
	fuzzBigIntPoly(p)
	str := fmt.Sprint(p)
	if str != "0" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = newBigIntPoly(N, R)
	p.Set(*big.NewInt(2), *big.NewInt(3), N)
	fuzzBigIntPoly(p)
	str = fmt.Sprint(p)
	if str != "x^3 + 2" {
		t.Error(dumpBigIntPoly(p), str)
	}

	p = newBigIntPoly(N, R)
	p.Set(*big.NewInt(1), *big.NewInt(1), N)
	fuzzBigIntPoly(p)
	str = fmt.Sprint(p)
	if str != "x + 1" {
		t.Error(dumpBigIntPoly(p), str)
	}
}
