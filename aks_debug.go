package main

import "fmt"
import "log"
import "math/big"
import "os"
import "runtime"
import "runtime/pprof"

func testAKSWitnesses(ch chan int) {
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

	ch <- 0
}

func main() {
	runtime.GOMAXPROCS(1)

	f, err := os.Create("cpu.out")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	ch := make(chan int, 1)
	go testAKSWitnesses(ch)

	<-ch
}
