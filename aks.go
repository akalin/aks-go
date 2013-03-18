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

func main() {
	n := big.NewInt(175507)
	r := big.NewInt(337)
	a := big.NewInt(2)
	fmt.Printf("n = %v, r = %v, a = %v\n", n, r, a)
	isWitness := isAKSWitness(n, r, a)
	var isWitnessStr string
	if isWitness {
		isWitnessStr = "is"
	} else {
		isWitnessStr = "is not"
	}
	fmt.Printf("a %s an AKS witness for n\n", isWitnessStr)
}
