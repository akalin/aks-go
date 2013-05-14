package main

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
	phi.Lsh(big.NewInt(1), uint(k*_BIG_WORD_BITS))
	phi.Add(&phi, big.NewInt(1))

	s := uint(R * k * _BIG_WORD_BITS)
	for i := 0; i < 45; i++ {
		fmt.Printf("%d: multiplying...\n", i)
		phi.Mul(&phi, &phi)
		fmt.Printf("%d: multiplying done; shifting...\n", i)
		len := uint(phi.BitLen())
		if len > s {
			fmt.Printf("%d: shifting...\n", i)
			phi.Rsh(&phi, len-s)
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
	resultCh chan witnessResult) {
	for a := range numberCh {
		fmt.Printf("Testing %v...\n", a)
		isAKSWitness()
		fmt.Printf("Finished testing %v\n", a)
		resultCh <- witnessResult{a, false}
	}
}

func main() {
	runtime.GOMAXPROCS(1)

	f, err := os.Create("cpu.out")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	numberCh := make(chan *big.Int, 1)
	defer close(numberCh)
	resultCh := make(chan witnessResult, 1)
	go testAKSWitnesses(numberCh, resultCh)

	for i := 1; i < 10; {
		select {
		case result := <-resultCh:
			fmt.Printf("%v isWitness=%t\n",
				result.a, result.isWitness)
			if result.isWitness {
				return
			}
		default:
			var a big.Int
			numberCh <- &a
		}
	}
}
