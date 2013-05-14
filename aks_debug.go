package main

import "log"
import "math/big"
import "os"
import "runtime/pprof"

func testAKSWitnesses(ch chan int) {
	R := 16451
	var bits uint = 5 * 64

	var phi big.Int
	phi.Lsh(big.NewInt(1), bits)
	phi.Add(&phi, big.NewInt(1))

	s := uint(R) * bits
	for i := 0; i < 15; i++ {
		log.Printf("%d: multiplying...\n", i)
		phi.Mul(&phi, &phi)
		log.Printf("%d: multiplying done; shifting...\n", i)
		len := uint(phi.BitLen())
		if len > s {
			log.Printf("%d: shifting...\n", i)
			phi.Rsh(&phi, len-s)
			log.Printf("%d: shifting done.\n", i)
		} else {
			log.Printf("%d: not shifting\n", i)
		}
	}

	ch <- 0
}

func main() {
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
