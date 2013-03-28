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

	if n.Sign() < 0 {
		panic("negative n")
	}
	if n.Sign() == 0 {
		return
	}
	t := &big.Int{}
	t.Set(n)
	// TODO(akalin): Use a wheel.
	for d := two; ; d.Add(d, one) {
		var dSq big.Int
		dSq.Mul(d, d)
		if dSq.Cmp(n) > 0 {
			break
		}
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
				return
			}
		}
	}
	if t.Cmp(one) != 0 {
		factorFn(t, one)
	}
}
