package aks

import "math/big"
import "testing"
import "unsafe"

// bigIntAsMpn() should return a pointer to the first element of its
// given big.Int's bits, and the length or capacity of its given
// big.Int's bits.
func TestBigIntAsMpn(t *testing.T) {
	x := (1 << 31) - 1
	i := big.NewInt(int64(x))
	bits := i.Bits()
	n1, s1 := bigIntAsMpn(i, _LEN)
	if unsafe.Pointer(n1) != unsafe.Pointer(&bits[0]) {
		t.Error(n1, &bits[0])
		return
	}
	if s1 != Size(len(bits)) {
		t.Error(s1, len(bits))
		return
	}

	n2, s2 := bigIntAsMpn(i, _CAP)
	if unsafe.Pointer(n2) != unsafe.Pointer(&bits[0]) {
		t.Error(n2, &bits[0])
		return
	}
	if s2 != Size(cap(bits)) {
		t.Error(s2, cap(bits))
		return
	}
}

// mpnMul() should perform multiplication on its operands.
func TestMpnMul(t *testing.T) {
	limbs1 := [3]Limb{1, 2, 3}
	limbs2 := [2]Limb{4, 5}
	limbs3 := [5]Limb{}

	mpnMul(&limbs3[0], &limbs1[0], 3, &limbs2[0], 2)
	expectedLimbs3 := [5]Limb{4, 13, 22, 15, 0}
	if limbs3 != expectedLimbs3 {
		t.Error(limbs3, expectedLimbs3)
	}
}

// mpnSqr() should perform squaring on its operand.
func TestMpnSqr(t *testing.T) {
	limbs1 := [3]Limb{1, 2, 3}
	limbs2 := [6]Limb{}

	mpnSqr(&limbs2[0], &limbs1[0], 3)
	expectedLimbs2 := [6]Limb{1, 4, 10, 12, 9, 0}
	if limbs2 != expectedLimbs2 {
		t.Error(limbs2, expectedLimbs2)
	}
}

// mpnTdivQr() should perform division with remainder on its operands.
func TestMpnTdivQr(t *testing.T) {
	limbs1 := [4]Limb{10, 13, 22, 15}
	limbs2 := [3]Limb{1, 2, 3}
	limbs3 := [2]Limb{}
	limbs4 := [3]Limb{}
	mpnTdivQr(&limbs3[0], &limbs4[0], 0, &limbs1[0], 4, &limbs2[0], 3)
	expectedLimbs3 := [2]Limb{4, 5}
	if limbs3 != expectedLimbs3 {
		t.Error(limbs3, expectedLimbs3)
	}
	expectedLimbs4 := [3]Limb{6, 0}
	if limbs4 != expectedLimbs4 {
		t.Error(limbs4, expectedLimbs4)
	}
}
