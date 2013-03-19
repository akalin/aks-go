package main

import "fmt"
import "math/big"

// Returns whether (X + a)^n = X^n + a mod (n, X^r - 1).
func isAKSWitness(n, r, a *big.Int) bool {
	reduceAKS := func(p *IntPoly) {
		p.PowMod(p, r).Mod(p, n)
	}

	zero := big.NewInt(0)
	one := big.NewInt(1)
	lhs := NewIntPoly([][2]*big.Int{{a, zero}, {one, one}})
	lhs.GenPow(lhs, n, reduceAKS)

	rhs := NewIntPoly([][2]*big.Int{{a, zero}, {one, n}})
	reduceAKS(rhs)

	isWitness := !lhs.Eq(rhs)
	return isWitness
}

// Returns the first AKS witness of n with the parameters r and M, or
// nil if there isn't one.
func getFirstAKSWitness(n, r, M *big.Int) *big.Int {
	for a := big.NewInt(1); a.Cmp(M) < 0; a.Add(a, big.NewInt(1)) {
		fmt.Printf("Testing %v (M = %v)...\n", a, M)
		if isWitness := isAKSWitness(n, r, a); isWitness {
			return a
		}
	}
	return nil
}

func main() {
	n := big.NewInt(46633)
	r := big.NewInt(262)
	M := big.NewInt(257)
	fmt.Printf("n = %v, r = %v, M = %v\n", n, r, M)
	a := getFirstAKSWitness(n, r, M)
	if a != nil {
		fmt.Printf("n is composite with AKS witness %v\n", a)
	} else {
		fmt.Printf("n is prime\n")
	}
}
