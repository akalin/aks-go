package main

import "math/big"

// Assuming p is prime, calculates and returns Phi(p^k) quickly.
func CalculateEulerPhiPrimePower(p, k *big.Int) *big.Int {
	var pMinusOne, kMinusOne big.Int
	pMinusOne.Sub(p, big.NewInt(1))
	kMinusOne.Sub(k, big.NewInt(1))
	var phi big.Int
	phi.Exp(p, &kMinusOne, nil)
	phi.Mul(&phi, &pMinusOne)
	return &phi
}

// A FactorFunction takes a prime and its multiplicity and returns
// whether or not to continue trying to find more factors.
type FactorFunction func(p, m *big.Int) bool

// Does trial division to find factors of n and passes them to the
// given FactorFunction until it indicates otherwise.
func TrialDivide(n *big.Int, factorFn FactorFunction) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)
	four := big.NewInt(4)
	five := big.NewInt(5)
	six := big.NewInt(6)
	seven := big.NewInt(7)
	nine := big.NewInt(9)
	eleven := big.NewInt(11)
	twentyFive := big.NewInt(25)
	fortyNine := big.NewInt(49)

	if n.Sign() < 0 {
		panic("negative n")
	}
	if n.Sign() == 0 {
		return
	}

	t := &big.Int{}
	t.Set(n)
	// Factors out d from t as much as possible and calls factorFn
	// if d divides t.
	factorOut := func(d *big.Int) bool {
		var m big.Int
		for {
			var q, r big.Int
			q.QuoRem(t, d, &r)
			if r.Sign() != 0 {
				break
			}
			t = &q
			m.Add(&m, one)
		}
		if m.Sign() != 0 {
			if !factorFn(d, &m) {
				return false
			}
		}
		return true
	}

	// TODO(akalin): Compute floor(sqrt(n)) once and compare d to
	// that instead of squaring d and comparing that to n.

	// Try small primes first.
	if four.Cmp(n) <= 0 && !factorOut(two) {
		return
	}

	if three.Cmp(t) <= 0 && nine.Cmp(n) <= 0 && !factorOut(three) {
		return
	}

	if five.Cmp(t) <= 0 && twentyFive.Cmp(n) <= 0 && !factorOut(five) {
		return
	}

	if seven.Cmp(t) <= 0 && fortyNine.Cmp(n) <= 0 && !factorOut(seven) {
		return
	}

	// Then run through a mod-30 wheel, which cuts the number of
	// odd numbers to test roughly in half.
	mod30Wheel := []*big.Int{four, two, four, two, four, six, two, six}
	for i, d := 1, eleven; d.Cmp(t) <= 0; {
		var dSq big.Int
		dSq.Mul(d, d)
		if dSq.Cmp(n) > 0 {
			break
		}
		if !factorOut(d) {
			return
		}
		d.Add(d, mod30Wheel[i])
		i = (i + 1) % len(mod30Wheel)
	}
	if t.Cmp(one) != 0 {
		factorFn(t, one)
	}
}
