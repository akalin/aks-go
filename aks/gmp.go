package aks

/*
#cgo LDFLAGS: -lgmp
#include <gmp.h>
*/
import "C"
import "math/big"
import "unsafe"

type Limb C.mp_limb_t
type Size C.mp_size_t

type sizeType int

const (
	_LEN sizeType = iota
	_CAP          = iota
)

// Return a pointer to the first Limb and the number of Limbs in the
// given big.Int (depending on the sizeType passed in). Assumes that
// sizeof(big.Word) == sizeof(Limb).
func bigIntAsMpn(i *big.Int, sizeType sizeType) (*Limb, Size) {
	bits := i.Bits()
	var size Size
	switch sizeType {
	case _LEN:
		size = Size(len(bits))
	case _CAP:
		size = Size(cap(bits))
	}
	if size == 0 {
		panic("empty big.Int")
	}
	allBits := bits[0:cap(bits)]
	return (*Limb)(unsafe.Pointer(&allBits[0])), size
}

// Multiply {s1p, s1n} and {s2p, s2n}, and write the (s1n+s2n)-limb
// result to rp. Return the most significant limb of the result.
//
// The destination has to have space for s1n + s2n limbs, even if the
// product's most significant limb is zero. No overlap is permitted
// between the destination and either source.
//
// This function requires that s1n is greater than or equal to s2n.
func mpnMul(rp, s1p *Limb, s1n Size, s2p *Limb, s2n Size) {
	C.mpn_mul(
		(*C.mp_limb_t)(rp),
		(*C.mp_limb_t)(s1p),
		C.mp_size_t(s1n),
		(*C.mp_limb_t)(s2p),
		C.mp_size_t(s2n))
}

// Compute the square of {s1p, n} and write the 2*n-limb result to rp.
//
// The destination has to have space for 2*n limbs, even if the
// result's most significant limb is zero. No overlap is permitted
// between the destination and the source.
func mpnSqr(rp, s1p *Limb, n Size) {
	C.mpn_sqr(
		(*C.mp_limb_t)(rp),
		(*C.mp_limb_t)(s1p),
		C.mp_size_t(n))
}
