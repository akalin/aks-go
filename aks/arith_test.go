package aks

import "math/big"
import "testing"

func expSmall(x, y int64) int64 {
	var z big.Int
	z.Exp(big.NewInt(x), big.NewInt(y), nil)
	return z.Int64()
}

func floorRootSmall(x, y int64) int64 {
	return floorRoot(big.NewInt(x), big.NewInt(y)).Int64()
}

// floorRoot(x^y, y) should always yield x.
func TestFloorRootExactPowers(t *testing.T) {
	for i := int64(0); i < 16; i++ {
		for j := int64(1); j < 16; j++ {
			k := floorRootSmall(expSmall(i, j), j)
			if k != i {
				t.Error(i, j, k)
			}
		}
	}
}

// floorRoot(x^y + 1, y) should yield x for x >= 1 and y >= 2.
func TestFloorRootSlightlyOverExactPower(t *testing.T) {
	for i := int64(1); i < 16; i++ {
		for j := int64(2); j < 16; j++ {
			k := floorRootSmall(expSmall(i, j)+1, j)
			if k != i {
				t.Error(i, j, k)
			}
		}
	}
}

// floorRoot((x + 1)^y - 1, y) should yield x for x >= 1 and y >= 2.
func TestFloorRootSlightlyUnderExactPower(t *testing.T) {
	for i := int64(1); i < 16; i++ {
		for j := int64(2); j < 16; j++ {
			k := floorRootSmall(expSmall(i+1, j)-1, j)
			if k != i {
				t.Error(i, j, k)
			}
		}
	}
}

// floorRoot((x^y + (x + 1)^y) / 2, y) should yield x for x >= 1 and y
// >= 2.
func TestFloorRootMidwayBetweenExactPowers(t *testing.T) {
	for i := int64(1); i < 16; i++ {
		for j := int64(2); j < 16; j++ {
			m := (expSmall(i, j) + expSmall(i+1, j)) / 2
			k := floorRootSmall(m, j)
			if k != i {
				t.Error(i, j, k)
			}
		}
	}
}

// Phi(p) should return p-1 for prime p.
func TestCalculateEulerPhiPrime(t *testing.T) {
	one := big.NewInt(1)

	phi := calculateEulerPhiPrimePower(big.NewInt(2), one)
	if phi.Cmp(one) != 0 {
		t.Error(phi)
	}

	phi = calculateEulerPhiPrimePower(big.NewInt(3), one)
	if phi.Cmp(big.NewInt(2)) != 0 {
		t.Error(phi)
	}

	phi = calculateEulerPhiPrimePower(big.NewInt(103), one)
	if phi.Cmp(big.NewInt(102)) != 0 {
		t.Error(phi)
	}
}

// Phi(p^k) should return p^(k-1)*(p-1) for prime p.
func TestCalculateEulerPhiPrimePower(t *testing.T) {
	phi := calculateEulerPhiPrimePower(big.NewInt(3), big.NewInt(5))
	if phi.Cmp(big.NewInt(162)) != 0 {
		t.Error(phi)
	}
}

// Converts a list of int64 pairs to a list of *big.Int pairs.
func makeFactors(int64Factors [][2]int64) [][2]*big.Int {
	factors := make([][2]*big.Int, len(int64Factors))
	for i, int64Factor := range int64Factors {
		factors[i][0] = big.NewInt(int64Factor[0])
		factors[i][1] = big.NewInt(int64Factor[1])
	}
	return factors
}

// Returns a factorFunction which compares its given factors to each
// successive element in the given list of factors.
func makeExpectingFactorFunction(
	n int64,
	int64Factors [][2]int64,
	comparedFactors *int,
	t *testing.T) factorFunction {
	expectedFactors := makeFactors(int64Factors)
	*comparedFactors = 0
	return func(p, m *big.Int) bool {
		if *comparedFactors >= len(expectedFactors) {
			t.Error(n, len(expectedFactors))
			return false
		}
		expectedP := expectedFactors[*comparedFactors][0]
		if p.Cmp(expectedP) != 0 {
			t.Error(n, p, expectedP)
			return false
		}
		expectedM := expectedFactors[*comparedFactors][1]
		if m.Cmp(expectedM) != 0 {
			t.Error(n, m, expectedM)
			return false
		}
		*comparedFactors++
		return true
	}
}

// Tests that trialDivide run with the given number gives the expected
// list of factors.
func testTrialDivide(n int64, expectedFactors [][2]int64, t *testing.T) {
	comparedFactors := 0
	trialDivide(
		big.NewInt(n),
		makeExpectingFactorFunction(
			n, expectedFactors, &comparedFactors, t),
		nil)
	if comparedFactors != len(expectedFactors) {
		t.Error(n, comparedFactors, len(expectedFactors))
	}
}

// Test trialDivide with small numbers.
func TestTrialDivideSmall(t *testing.T) {
	testTrialDivide(0, [][2]int64{}, t)
	testTrialDivide(1, [][2]int64{}, t)
	testTrialDivide(2, [][2]int64{{2, 1}}, t)
	testTrialDivide(3, [][2]int64{{3, 1}}, t)
	testTrialDivide(4, [][2]int64{{2, 2}}, t)
	testTrialDivide(5, [][2]int64{{5, 1}}, t)
	testTrialDivide(6, [][2]int64{{2, 1}, {3, 1}}, t)
	testTrialDivide(7, [][2]int64{{7, 1}}, t)
	testTrialDivide(8, [][2]int64{{2, 3}}, t)
	testTrialDivide(9, [][2]int64{{3, 2}}, t)
	testTrialDivide(10, [][2]int64{{2, 1}, {5, 1}}, t)
}

// Test trialDivide with some larger numbers.
func TestTrialDivideLarge(t *testing.T) {
	testTrialDivide(100, [][2]int64{{2, 2}, {5, 2}}, t)
	testTrialDivide(101, [][2]int64{{101, 1}}, t)
	testTrialDivide(1961, [][2]int64{{37, 1}, {53, 1}}, t)
}

// Make sure trialDivide respects the return value of its
// factorFunction.
func TestTrialDividePartial(t *testing.T) {
	var n int64 = 100
	expectedFactors := [][2]int64{{2, 2}}
	comparedFactors := 0
	expectingFactorFunction :=
		makeExpectingFactorFunction(
			n, expectedFactors, &comparedFactors, t)
	partialFactorFunction := func(p, m *big.Int) bool {
		if comparedFactors >= 1 {
			return false
		}
		return expectingFactorFunction(p, m)
	}
	trialDivide(big.NewInt(n), partialFactorFunction, nil)
	if comparedFactors != len(expectedFactors) {
		t.Error(comparedFactors, len(expectedFactors))
	}
}

func calculateMultiplicativeOrderPrimePowerSmall(a, p, k int64) int64 {
	return calculateMultiplicativeOrderPrimePower(
		big.NewInt(a), big.NewInt(p), big.NewInt(k)).Int64()
}

// Check calculateMultiplicativeOrderPrimePower() with some small test
// cases.
func TestCalculateMultiplicativeOrderPrimePower(t *testing.T) {
	e := calculateMultiplicativeOrderPrimePowerSmall(4, 7, 1)
	if e != 3 {
		t.Error(e)
	}

	e = calculateMultiplicativeOrderPrimePowerSmall(3, 2, 10)
	if e != 256 {
		t.Error(e)
	}
}

func calculateMultiplicativeOrderSmall(a, n int64) int64 {
	return calculateMultiplicativeOrder(
		big.NewInt(a), big.NewInt(n)).Int64()
}

// Check calculateMultiplicativeOrder() with a small test case.
func TestCalculateMultiplicativeOrder(t *testing.T) {
	e := calculateMultiplicativeOrderSmall(3, 25600)
	if e != 1280 {
		t.Error(e)
	}
}

// Check calculateEulerPhi() with a small test case.
func TestCalculateEulerPhi(t *testing.T) {
	// 3888 = 2^4 * 3^5.
	phi := calculateEulerPhi(big.NewInt(3888))
	// phi(3888) = phi(2^4) * phi(3^5) = 2^3 * 3^4 * 2 = 6^4 = 1296.
	if phi.Cmp(big.NewInt(1296)) != 0 {
		t.Error(phi)
	}
}
