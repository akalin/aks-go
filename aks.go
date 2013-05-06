package main

import "flag"
import "fmt"
import "log"
import "math/big"
import "os"
import "runtime"

// Returns whether (X + a)^n = X^n + a mod (n, X^r - 1). tmp1, tmp2,
// and tmp3 must be BigIntPoly objects constructed with N, R = n, r,
// and they must not alias each other.
func isAKSWitness(n, a big.Int, tmp1, tmp2, tmp3 *BigIntPoly) bool {
	// Left-hand side: (X + a)^n mod (n, X^r - 1).
	tmp1.Set(a, *big.NewInt(1), n)
	tmp1.Pow(n, tmp2, tmp3)

	// Right-hand side: (X^n + a) mod (n, X^r - 1).
	tmp2.Set(a, n, n)

	isWitness := !tmp1.Eq(tmp2)
	return isWitness
}

// Returns the first AKS witness of n with the parameters r and M, or
// nil if there isn't one.
func getFirstAKSWitness(n, r, M *big.Int, logger *log.Logger) *big.Int {
	tmp1 := NewBigIntPoly(*n, *r)
	tmp2 := NewBigIntPoly(*n, *r)
	tmp3 := NewBigIntPoly(*n, *r)

	for a := big.NewInt(1); a.Cmp(M) < 0; a.Add(a, big.NewInt(1)) {
		logger.Printf("Testing %v (M = %v)...\n", a, M)
		isWitness := isAKSWitness(*n, *a, tmp1, tmp2, tmp3)
		if isWitness {
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
	resultCh chan witnessResult,
	logger *log.Logger) {
	tmp1 := NewBigIntPoly(*n, *r)
	tmp2 := NewBigIntPoly(*n, *r)
	tmp3 := NewBigIntPoly(*n, *r)

	for a := range numberCh {
		logger.Printf("Testing %v...\n", a)
		isWitness := isAKSWitness(*n, *a, tmp1, tmp2, tmp3)
		logger.Printf("Finished testing %v (isWitness=%t)\n",
			a, isWitness)
		resultCh <- witnessResult{a, isWitness}
	}
}

// Returns an AKS witness of n with the parameters r and M, or nil if
// there isn't one. Tests up to maxOutstanding numbers at once.
func getAKSWitness(
	n, r, M *big.Int,
	maxOutstanding int,
	logger *log.Logger) *big.Int {
	numberCh := make(chan *big.Int, maxOutstanding)
	defer close(numberCh)
	resultCh := make(chan witnessResult, maxOutstanding)
	for i := 0; i < maxOutstanding; i++ {
		go testAKSWitnesses(n, r, numberCh, resultCh, logger)
	}

	// Send off all numbers for testing, draining any results that
	// come in while we're doing so.
	tested := big.NewInt(1)
	for i := big.NewInt(1); i.Cmp(M) < 0; {
		select {
		case result := <-resultCh:
			tested.Add(tested, big.NewInt(1))
			logger.Printf("%v isWitness=%t\n",
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
		logger.Printf("%v isWitness=%t\n", result.a, result.isWitness)
		if result.isWitness {
			return result.a
		}
	}

	return nil
}

// Returns an upper bound for r such that o_r(n) > ceil(lg(n))^2 that
// is polylog in n.
func calculateAKSModulusUpperBound(n *big.Int) *big.Int {
	two := big.NewInt(2)
	three := big.NewInt(3)
	five := big.NewInt(5)
	eight := big.NewInt(8)

	// Calculate max(ceil(lg(n))^5, 3).
	ceilLgN := big.NewInt(int64(n.BitLen()))
	rUpperBound := &big.Int{}
	rUpperBound.Exp(ceilLgN, five, nil)
	rUpperBound = Max(rUpperBound, three)

	var nMod8 big.Int
	nMod8.Mod(n, eight)
	if (nMod8.Cmp(three) == 0) || (nMod8.Cmp(five) == 0) {
		// Calculate 8*ceil(lg(n))^2.
		var rUpperBound2 big.Int
		rUpperBound2.Exp(ceilLgN, two, nil)
		rUpperBound2.Mul(&rUpperBound2, eight)
		rUpperBound = Min(rUpperBound, &rUpperBound2)
	}
	return rUpperBound
}

// Returns the least r such that o_r(n) > ceil(lg(n))^2 >= ceil(lg(n)^2).
func calculateAKSModulus(n *big.Int) *big.Int {
	one := big.NewInt(1)
	two := big.NewInt(2)

	ceilLgNSq := big.NewInt(int64(n.BitLen()))
	ceilLgNSq.Mul(ceilLgNSq, ceilLgNSq)
	var r big.Int
	r.Add(ceilLgNSq, two)
	rUpperBound := calculateAKSModulusUpperBound(n)
	for ; r.Cmp(rUpperBound) < 0; r.Add(&r, one) {
		var gcd big.Int
		gcd.GCD(nil, nil, n, &r)
		if gcd.Cmp(one) != 0 {
			continue
		}
		o := CalculateMultiplicativeOrder(n, &r)
		if o.Cmp(ceilLgNSq) > 0 {
			return &r
		}
	}

	panic("Could not find AKS modulus")
}

// Returns floor(sqrt(Phi(r))) * ceil(lg(n)) + 1 > floor(sqrt(Phi(r))) * lg(n).
func calculateAKSUpperBound(n, r *big.Int) *big.Int {
	one := big.NewInt(1)
	two := big.NewInt(2)

	M := CalculateEulerPhi(r)
	M = FloorRoot(M, two)
	M.Mul(M, big.NewInt(int64(n.BitLen())))
	M.Add(M, one)
	return M
}

// Returns the first factor of n less than M.
func getFirstFactorBelow(n, M *big.Int) *big.Int {
	var factor *big.Int
	var mMinusOne big.Int
	mMinusOne.Sub(M, big.NewInt(1))
	TrialDivide(n, func(q, e *big.Int) bool {
		if q.Cmp(M) < 0 && q.Cmp(n) < 0 {
			factor = q
		}
		return false
	}, &mMinusOne)
	return factor
}

func main() {
	jobs := flag.Int(
		"j", runtime.NumCPU(), "how many processing jobs to spawn")

	flag.Parse()

	runtime.GOMAXPROCS(*jobs)

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "%s [options] [number]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	var n big.Int
	_, parsed := n.SetString(flag.Arg(0), 10)
	if !parsed {
		fmt.Fprintf(os.Stderr, "could not parse %s\n", flag.Arg(0))
		os.Exit(-1)
	}

	two := big.NewInt(2)

	if n.Cmp(two) < 0 {
		fmt.Fprintf(os.Stderr, "n must be >= 2\n")
		os.Exit(-1)
	}

	r := calculateAKSModulus(&n)
	M := calculateAKSUpperBound(&n, r)
	fmt.Printf("n = %v, r = %v, M = %v\n", &n, r, M)
	factor := getFirstFactorBelow(&n, M)
	if factor != nil {
		fmt.Printf("n has factor %v\n", factor)
		return
	}

	fmt.Printf("n has no factor less than %v\n", M)
	sqrtN := FloorRoot(&n, two)
	if M.Cmp(sqrtN) > 0 {
		fmt.Printf("%v is greater than sqrt(%v), so %v is prime\n",
			M, &n, &n)
		return
	}

	a := getAKSWitness(&n, r, M, *jobs, log.New(os.Stderr, "", 0))
	if a != nil {
		fmt.Printf("n is composite with AKS witness %v\n", a)
	} else {
		fmt.Printf("n is prime\n")
	}
}
