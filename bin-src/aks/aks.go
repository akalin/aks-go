package main

import "github.com/akalin/aks-go/aks"
import "flag"
import "fmt"
import "log"
import "math/big"
import "os"
import "runtime"
import "runtime/pprof"

func main() {
	jobs := flag.Int(
		"j", runtime.NumCPU(), "how many processing jobs to spawn")
	endStr := flag.String(
		"end", "", "the upper bound to use (defaults to M)")
	cpuProfilePath :=
		flag.String("cpuprofile", "",
			"Write a CPU profile to the specified file "+
				"before exiting.")

	flag.Parse()

	runtime.GOMAXPROCS(*jobs)

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

	r := aks.CalculateAKSModulus(&n)
	M := aks.CalculateAKSUpperBound(&n, r)

	if end.Sign() <= 0 {
		end.Set(M)
	}
	fmt.Printf("n = %v, r = %v, M = %v, end = %v\n", &n, r, M, &end)
	factor := aks.GetFirstFactorBelow(&n, M)
	if factor != nil {
		fmt.Printf("n has factor %v\n", factor)
		return
	}

	fmt.Printf("n has no factor less than %v\n", M)
	// M^2 > N iff M > floor(sqrt(N)).
	var mSq big.Int
	mSq.Mul(M, M)
	if mSq.Cmp(&n) > 0 {
		fmt.Printf("%v is greater than sqrt(%v), so %v is prime\n",
			M, &n, &n)
		return
	}

	a := aks.GetAKSWitness(&n, r, &end, *jobs, log.New(os.Stderr, "", 0))
	if a != nil {
		fmt.Printf("n is composite with AKS witness %v\n", a)
	} else if end.Cmp(M) < 0 {
		fmt.Printf("n has no AKS witnesses < %v\n", &end)
	} else {
		fmt.Printf("n is prime\n")
	}
}
