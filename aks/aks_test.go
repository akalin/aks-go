package aks

import "io/ioutil"
import "log"
import "math/big"
import "runtime"
import "testing"

// The number of rounds to use for big.Int.ProbablyPrime().
const _NUM_PROBABLY_PRIME_ROUNDS = 10

// Returns the first prime with the given number of decimal digits.
func getFirstPrimeWithDigits(numDigits int64) *big.Int {
	one := big.NewInt(1)
	n := big.NewInt(10)
	n.Exp(n, big.NewInt(numDigits), nil)
	for !n.ProbablyPrime(_NUM_PROBABLY_PRIME_ROUNDS) {
		n.Add(n, one)
	}
	return n
}

// Benchmark isAKSWitness for the first prime number of the given
// number of decimal digits.
func runIsAKSWitnessBenchmark(b *testing.B, numDigits int64) {
	b.StopTimer()
	n := getFirstPrimeWithDigits(numDigits)
	r := CalculateAKSModulus(n)
	// Any a > 1 suffices.
	a := big.NewInt(2)

	tmp1 := newBigIntPoly(*n, *r)
	tmp2 := newBigIntPoly(*n, *r)
	tmp3 := newBigIntPoly(*n, *r)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		isAKSWitness(*n, *a, tmp1, tmp2, tmp3)
	}
}

// Benchmark isAKSWitness for values of n of varying digit sizes.

func BenchmarkIsAKSWitness3Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 3)
}

func BenchmarkIsAKSWitness4Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 4)
}

func BenchmarkIsAKSWitness5Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 5)
}

func BenchmarkIsAKSWitness6Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 6)
}

func BenchmarkIsAKSWitness7Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 7)
}

func BenchmarkIsAKSWitness8Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 8)
}

func BenchmarkIsAKSWitnessMax32(b *testing.B) {
	b.StopTimer()
	n := big.NewInt(4294967291)
	r := CalculateAKSModulus(n)
	// Any a > 1 suffices.
	a := big.NewInt(2)

	tmp1 := newBigIntPoly(*n, *r)
	tmp2 := newBigIntPoly(*n, *r)
	tmp3 := newBigIntPoly(*n, *r)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		isAKSWitness(*n, *a, tmp1, tmp2, tmp3)
	}
}

func BenchmarkIsAKSWitness9Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 9)
}

func BenchmarkIsAKSWitness10Digits(b *testing.B) {
	runIsAKSWitnessBenchmark(b, 10)
}

var nullLogger *log.Logger = log.New(ioutil.Discard, "", 0)

// Benchmark getFirstAKSWitness for the first prime number of the
// given number of decimal digits.
func runGetFirstAKSWitnessBenchmark(b *testing.B, numDigits int64) {
	b.StopTimer()
	n := getFirstPrimeWithDigits(numDigits)
	r := CalculateAKSModulus(n)
	M := big.NewInt(10)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getFirstAKSWitness(n, r, M, nullLogger)
	}
}

// Benchmark getFirstAKSWitness for values of n of varying digit sizes.

func BenchmarkGetFirstAKSWitness3Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 3)
}

func BenchmarkGetFirstAKSWitness4Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 4)
}

func BenchmarkGetFirstAKSWitness5Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 5)
}

func BenchmarkGetFirstAKSWitness6Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 6)
}

func BenchmarkGetFirstAKSWitness7Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 7)
}

func BenchmarkGetFirstAKSWitness8Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 8)
}

func BenchmarkGetFirstAKSWitness11Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 11)
}

func BenchmarkGetFirstAKSWitness12Digits(b *testing.B) {
	runGetFirstAKSWitnessBenchmark(b, 12)
}

// Benchmark getAKSWitness for the first prime number of the given
// number of decimal digits.
func runGetAKSWitnessBenchmark(b *testing.B, numDigits int64) {
	b.StopTimer()
	n := getFirstPrimeWithDigits(numDigits)
	r := CalculateAKSModulus(n)
	M := big.NewInt(10)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetAKSWitness(n, r, M, runtime.GOMAXPROCS(0), nullLogger)
	}
}

// Benchmark getFirstAKSWitness for values of n of varying digit sizes.

func BenchmarkGetAKSWitness3Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 3)
}

func BenchmarkGetAKSWitness4Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 4)
}

func BenchmarkGetAKSWitness5Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 5)
}

func BenchmarkGetAKSWitness6Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 6)
}

func BenchmarkGetAKSWitness7Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 7)
}

func BenchmarkGetAKSWitness8Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 8)
}

func BenchmarkGetAKSWitness11Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 11)
}

func BenchmarkGetAKSWitness12Digits(b *testing.B) {
	runGetAKSWitnessBenchmark(b, 12)
}
