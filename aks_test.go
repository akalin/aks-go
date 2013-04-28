package main

import "io/ioutil"
import "log"
import "math/big"
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
	r := calculateAKSModulus(n)
	// Any a > 1 suffices.
	a := big.NewInt(2)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		isAKSWitness(n, r, a)
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

var nullLogger *log.Logger = log.New(ioutil.Discard, "", 0)

// Benchmark getFirstAKSWitness for the first prime number of the
// given number of decimal digits.
func runGetFirstAKSWitnessBenchmark(b *testing.B, numDigits int64) {
	b.StopTimer()
	n := getFirstPrimeWithDigits(numDigits)
	r := calculateAKSModulus(n)
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
