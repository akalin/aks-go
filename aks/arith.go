package aks

import "math/big"

// Returns the smaller of x and y. No copies are made, so the returned
// pointer is either x or y.
func min(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return x
	}
	return y
}

// Returns the larger of x and y. No copies are made, so the returned
// pointer is either x or y.
func max(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return x
	}
	return y
}

// Returns the greatest number y such that y^k <= x. x must be
// non-negative and k must be positive.
func floorRoot(x, k *big.Int) *big.Int {
	if x.Sign() < 0 {
		panic("negative radicand")
	}
	if k.Sign() <= 0 {
		panic("non-negative index")
	}
	if x.Sign() == 0 {
		return &big.Int{}
	}
	one := big.NewInt(1)
	var kMinusOne big.Int
	kMinusOne.Sub(k, one)

	// Calculate p = ceil((floor(lg(x)) + 1)/k).
	var p, r big.Int
	p.DivMod(big.NewInt(int64(x.BitLen())), k, &r)
	if r.Sign() > 0 {
		p.Add(&p, one)
	}

	y := &big.Int{}
	y.Exp(big.NewInt(2), &p, nil)
	for y.Cmp(one) > 0 {
		// Calculate z = floor(((k-1)y + floor(x/y^{k-1}))/k).
		var z1 big.Int
		z1.Mul(&kMinusOne, y)

		var z2 big.Int
		var yPowKMinusOne big.Int
		yPowKMinusOne.Exp(y, &kMinusOne, nil)
		z2.Div(x, &yPowKMinusOne)

		var z big.Int
		z.Add(&z1, &z2)
		z.Div(&z, k)

		if z.Cmp(y) >= 0 {
			return y
		}
		y = &z
	}
	return one
}

// Assuming p is prime, calculates and returns Phi(p^k) quickly.
func calculateEulerPhiPrimePower(p, k *big.Int) *big.Int {
	var pMinusOne, kMinusOne big.Int
	pMinusOne.Sub(p, big.NewInt(1))
	kMinusOne.Sub(k, big.NewInt(1))
	var phi big.Int
	phi.Exp(p, &kMinusOne, nil)
	phi.Mul(&phi, &pMinusOne)
	return &phi
}

// A factorFunction takes a prime and its multiplicity and returns
// whether or not to continue trying to find more factors.
type factorFunction func(p, m *big.Int) bool

// Does trial division to find factors of n and passes them to the
// given factorFunction until it indicates otherwise. If upperBound is
// not nil, only factors less than or equal to it will be tried.
func trialDivide(n *big.Int, factorFn factorFunction, upperBound *big.Int) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)
	four := big.NewInt(4)
	five := big.NewInt(5)
	six := big.NewInt(6)
	seven := big.NewInt(7)
	eleven := big.NewInt(11)

	if n.Sign() < 0 {
		panic("negative n")
	}
	if n.Sign() == 0 {
		return
	}

	if upperBound == nil {
		upperBound = floorRoot(n, two)
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
			upperBound = min(upperBound, t)
			m.Add(&m, one)
		}
		if m.Sign() != 0 {
			if !factorFn(d, &m) {
				return false
			}
		}
		return true
	}

	// Try small primes first.
	if two.Cmp(upperBound) <= 0 && !factorOut(two) {
		return
	}

	if three.Cmp(upperBound) <= 0 && !factorOut(three) {
		return
	}

	if five.Cmp(upperBound) <= 0 && !factorOut(five) {
		return
	}

	if seven.Cmp(upperBound) <= 0 && !factorOut(seven) {
		return
	}

	// Then run through a mod-30 wheel, which cuts the number of
	// odd numbers to test roughly in half.
	mod30Wheel := []*big.Int{four, two, four, two, four, six, two, six}
	for i, d := 1, eleven; d.Cmp(upperBound) <= 0; {
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

// Assuming that p is prime and a and p^k are coprime, returns the
// smallest power e of a such that a^e = 1 (mod p^k).
func calculateMultiplicativeOrderPrimePower(a, p, k *big.Int) *big.Int {
	var n big.Int
	n.Exp(p, k, nil)
	t := calculateEulerPhiPrimePower(p, k)

	o := big.NewInt(1)
	one := big.NewInt(1)
	processPrimeFactor := func(q, e *big.Int) bool {
		// Calculate x = a^(t/q^e) (mod n).
		var x big.Int
		x.Exp(q, e, nil)
		x.Div(t, &x)
		x.Exp(a, &x, &n)
		for x.Cmp(one) != 0 {
			o.Mul(o, q)
			x.Exp(&x, q, &n)
		}
		return true
	}

	if k.Cmp(one) > 0 {
		var kMinusOne big.Int
		kMinusOne.Sub(k, one)
		processPrimeFactor(p, &kMinusOne)
	}

	var pMinusOne big.Int
	pMinusOne.Sub(p, one)
	trialDivide(&pMinusOne, processPrimeFactor, nil)

	return o
}

// Assuming that a and n are coprime, returns the smallest power e of
// a such that a^e = 1 (mod n).
func calculateMultiplicativeOrder(a, n *big.Int) *big.Int {
	o := big.NewInt(1)
	trialDivide(n, func(q, e *big.Int) bool {
		oq := calculateMultiplicativeOrderPrimePower(a, q, e)
		// Set o to lcm(o, oq).
		var gcd big.Int
		gcd.GCD(nil, nil, o, oq)
		o.Div(o, &gcd)
		o.Mul(o, oq)
		return true
	}, nil)
	return o
}

// Calculate Phi(n) by factorizing it.
func calculateEulerPhi(n *big.Int) *big.Int {
	phi := big.NewInt(1)
	trialDivide(n, func(q, e *big.Int) bool {
		phi.Mul(phi, calculateEulerPhiPrimePower(q, e))
		return true
	}, nil)
	return phi
}
