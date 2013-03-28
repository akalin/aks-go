package main

import "math/big"

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

	if four.Cmp(n) <= 0 && !factorOut(two) {
		return
	}

	// TODO(akalin): Use a mod-30 wheel.
	for d := three; d.Cmp(t) <= 0; d.Add(d, two) {
		var dSq big.Int
		dSq.Mul(d, d)
		if dSq.Cmp(n) > 0 {
			break
		}
		if !factorOut(d) {
			return
		}
	}
	if t.Cmp(one) != 0 {
		factorFn(t, one)
	}
}
