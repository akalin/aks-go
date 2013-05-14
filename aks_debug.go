package main

import "flag"
import "fmt"
import "log"
import "math/big"
import "os"
import "runtime"
import "runtime/pprof"

const (
	// Compute the size of a big.Word in bits.
	_m             = ^big.Word(0)
	_logS          = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S             = 1 << _logS
	_BIG_WORD_BITS = _S << 3
)

func isAKSWitness() {
	var n big.Int
	_, parsed := n.SetString("332315159569814711702351072539787810327", 10)
	if !parsed {
		panic("could not parse")
	}

	R := 16451

	var maxCoefficient big.Int
	maxCoefficient.Sub(&n, big.NewInt(1))
	maxCoefficient.Mul(&maxCoefficient, &maxCoefficient)
	maxCoefficient.Mul(&maxCoefficient, big.NewInt(int64(R)))

	k := len(maxCoefficient.Bits())

	var phi big.Int
	phi.Lsh(big.NewInt(1), uint(k * _BIG_WORD_BITS))
	phi.Add(&phi, big.NewInt(1))

	s := uint(R * k * _BIG_WORD_BITS)
	for i := 0; ; i++ {
		fmt.Printf("%d: multiplying...\n", i)
		phi.Mul(&phi, &phi)
		fmt.Printf("%d: multiplying done; shifting...\n", i)
		len := uint(phi.BitLen())
		if len > s {
			fmt.Printf("%d: shifting...\n", i)
			phi.Rsh(&phi, len - s)
			fmt.Printf("%d: shifting done.\n", i)
		} else {
			fmt.Printf("%d: not shifting\n", i)
		}
	}
}

// Holds the result of an AKS witness test.
type witnessResult struct {
	a         *big.Int
	isWitness bool
}

// Tests all numbers received on numberCh if they are witnesses of n
// with parameter r. Sends the results to resultCh.
func testAKSWitnesses(
	numberCh chan *big.Int,
	resultCh chan witnessResult,
	logger *log.Logger) {
	for a := range numberCh {
		logger.Printf("Testing %v...\n", a)
		isAKSWitness()
		logger.Printf("Finished testing %v\n", a)
		resultCh <- witnessResult{a, false}
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
		go testAKSWitnesses(numberCh, resultCh, logger)
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
	endStr := flag.String(
		"end", "", "the upper bound to use (defaults to M)")
	cpuProfilePath :=
		flag.String("cpuprofile", "",
			"Write a CPU profile to the specified file "+
				"before exiting.")

	flag.Parse()

	runtime.GOMAXPROCS(1)

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "%s [options] [number]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if len(*cpuProfilePath) > 0 {
		f, err := os.Create(*cpuProfilePath)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var end big.Int
	if len(*endStr) > 0 {
		_, parsed := end.SetString(*endStr, 10)
		if !parsed {
			fmt.Fprintf(os.Stderr, "could not parse %s\n", *endStr)
			os.Exit(-1)
		}
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

	if end.Sign() <= 0 {
		end.Set(M)
	}
	fmt.Printf("n = %v, r = %v, M = %v, end = %v\n", &n, r, M, &end)
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

	a := getAKSWitness(&n, r, &end, 1, log.New(os.Stderr, "", 0))
	if a != nil {
		fmt.Printf("n is composite with AKS witness %v\n", a)
	} else if end.Cmp(M) < 0 {
		fmt.Printf("n has no AKS witnesses < %v\n", &end)
	} else {
		fmt.Printf("n is prime\n")
	}
}
