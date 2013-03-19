package main

import "fmt"
import "math/big"
import "os"
import "runtime"

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

// Holds the result of an AKS witness test.
type witnessResult struct {
	a         *big.Int
	isWitness bool
}

// Tests all numbers received on numberCh if they are witnesses of n
// with parameter r. Sends the results to resultCh.
func testAKSWitnesses(
	n, r *big.Int,
	numberCh chan *big.Int,
	resultCh chan witnessResult) {
	for a := range numberCh {
		fmt.Printf("Testing %v...\n", a)
		isWitness := isAKSWitness(n, r, a)
		fmt.Printf("Finished testing %v (isWitness=%t)\n",
			a, isWitness)
		resultCh <- witnessResult{a, isWitness}
	}
}

// Returns an AKS witness of n with the parameters r and M, or nil if
// there isn't one. Tests up to maxOutstanding numbers at once.
func getAKSWitness(n, r, M *big.Int, maxOutstanding int) *big.Int {
	numberCh := make(chan *big.Int, maxOutstanding)
	defer close(numberCh)
	resultCh := make(chan witnessResult, maxOutstanding)
	for i := 0; i < maxOutstanding; i++ {
		go testAKSWitnesses(n, r, numberCh, resultCh)
	}

	// Send off all numbers for testing, draining any results that
	// come in while we're doing so.
	tested := big.NewInt(1)
	for i := big.NewInt(1); i.Cmp(M) < 0; {
		select {
		case result := <-resultCh:
			tested.Add(tested, big.NewInt(1))
			fmt.Printf("%v isWitness=%t\n",
				result.a, result.isWitness)
			if result.isWitness {
				return result.a
			}
		default:
			var a big.Int
			a.Set(i)
			numberCh <- &a
			i.Add(i, big.NewInt(1))
		}
	}

	// Drain any remaining results.
	for tested.Cmp(M) < 0 {
		result := <-resultCh
		tested.Add(tested, big.NewInt(1))
		fmt.Printf("%v isWitness=%t\n", result.a, result.isWitness)
		if result.isWitness {
			return result.a
		}
	}

	return nil
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s [number]\n", os.Args[0])
		os.Exit(-1)
	}

	var n big.Int
	_, parsed := n.SetString(os.Args[1], 10)
	if !parsed {
		fmt.Fprintf(os.Stderr, "could not parse %s\n", os.Args[1])
		os.Exit(-1)
	}
	if n.Cmp(big.NewInt(2)) < 0 {
		fmt.Fprintf(os.Stderr, "n must be >= 2\n")
		os.Exit(-1)
	}

	// TODO(akalin): Calculate AKS parameters properly.
	r := n
	M := n
	fmt.Printf("n = %v, r = %v, M = %v\n", &n, &r, &M)
	a := getAKSWitness(&n, &r, &M, numCPU)
	if a != nil {
		fmt.Printf("n is composite with AKS witness %v\n", a)
	} else {
		fmt.Printf("n is prime\n")
	}
}
