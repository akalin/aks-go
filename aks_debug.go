package main

import "fmt"
import "log"
import "math/big"
import "os"
import "runtime"
import "runtime/pprof"

func isAKSWitness() {
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
	<-numberCh

	R := 16451
	var bits uint = 5 * 64

	var phi big.Int
	phi.Lsh(big.NewInt(1), bits)
	phi.Add(&phi, big.NewInt(1))

	s := uint(R) * bits
	for i := 0; i < 15; i++ {
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

	resultCh <- witnessResult{nil, false}
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
	resultCh := make(chan witnessResult, 1)
	go testAKSWitnesses(numberCh, resultCh)

	var a big.Int
	numberCh <- &a
	<-resultCh
}
