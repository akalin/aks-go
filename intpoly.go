package main

import "math/big"

// An intMono represents the monomial coeff*X^deg.
// The zero value for an intMono represents the zero monomial.
type intMono struct {
	coeff big.Int
	deg   big.Int
}

// Returns whether m and n have the same coefficient and degree.
func (m *intMono) Eq(n *intMono) bool {
	return m.coeff.Cmp(&n.coeff) == 0 && m.deg.Cmp(&n.deg) == 0
}

// Sets m to a deep copy of n.
func (m *intMono) Set(n *intMono) {
	m.coeff.Set(&n.coeff)
	m.deg.Set(&n.deg)
}

// Sets m to the product of n and coeff*X^deg.
func (m *intMono) Mul(n *intMono, coeff, deg *big.Int) {
	m.coeff.Mul(&n.coeff, coeff)
	m.deg.Add(&n.deg, deg)
}

// Sets m to n with its coefficient reduced modulo k. If that ends up
// setting the coefficient to zero, then m becomes the zero monomial.
func (m *intMono) Mod(n *intMono, k *big.Int) {
	m.coeff.Mod(&n.coeff, k)
	if m.coeff.Sign() != 0 {
		m.deg.Set(&n.deg)
	}
}

// The container for a polynomial's terms.
type termList []intMono

// Returns a termList of the given size, possibly reusing term's
// space.
func (terms termList) make(n int) termList {
	if n <= cap(terms) {
		// Reuse the space.
		return terms[0:n]
	}
	// Otherwise, allocate a new array. No need to worry about
	// allocating extra space since we don't have to worry about
	// carries.
	return make(termList, n)
}

// An IntPoly represents the polynomial with the given non-zero terms
// in order of ascending degree.
// The zero value for an IntPoly represents the zero polynomial.
type IntPoly struct {
	terms termList
}

// Builds a new IntPoly from the given list of coefficient/degree
// pairs. Each coefficient must be non-zero, each degree must be
// non-negative, and the list must be in ascending order of degree.
func NewIntPoly(terms [][2]*big.Int) *IntPoly {
	var p IntPoly
	p.terms = p.terms.make(len(terms))
	for i, term := range terms {
		if term[0].Sign() == 0 {
			panic("zero coefficient")
		}
		if term[1].Sign() < 0 {
			panic("negative degree")
		}
		if i > 0 && term[1].Cmp(terms[i-1][1]) <= 0 {
			panic("non-increasing degree")
		}
		p.terms[i].coeff.Set(term[0])
		p.terms[i].deg.Set(term[1])
	}
	return &p
}

// Returns whether p and q have the same terms.
func (p *IntPoly) Eq(q *IntPoly) bool {
	if len(p.terms) != len(q.terms) {
		return false
	}
	for i, pTerm := range p.terms {
		if !pTerm.Eq(&q.terms[i]) {
			return false
		}
	}
	return true
}

// Sets p to a deep copy of q.
func (p *IntPoly) Set(q *IntPoly) *IntPoly {
	if p == q {
		return p
	}
	p.terms = p.terms.make(len(q.terms))
	for i, _ := range p.terms {
		p.terms[i].Set(&q.terms[i])
	}
	return p
}

// Sets p to the sum of q and r.
func (p *IntPoly) Add(q, r *IntPoly) *IntPoly {
	// Since we go left to right, reusing p's term array is okay
	// even if p aliases q or r (or both).
	terms := p.terms.make(len(q.terms) + len(r.terms))

	i, j, k := 0, 0, 0
	for ; j < len(q.terms) && k < len(r.terms); i++ {
		term := &terms[i]
		qTerm := &q.terms[j]
		rTerm := &r.terms[k]
		cmp := qTerm.deg.Cmp(&rTerm.deg)
		if cmp < 0 {
			term.Set(qTerm)
			j++
		} else if cmp > 0 {
			term.Set(rTerm)
			k++
		} else {
			term.coeff.Add(&qTerm.coeff, &rTerm.coeff)
			term.deg.Set(&qTerm.deg)
			j++
			k++
		}
	}

	if j < len(q.terms) {
		for ; j < len(q.terms); j++ {
			terms[i].Set(&q.terms[j])
			i++
		}
	} else if k < len(r.terms) {
		for ; k < len(r.terms); k++ {
			terms[i].Set(&r.terms[k])
			i++
		}
	}

	p.terms = terms[0:i]
	return p
}

// Sets p to the product of q and coeff*X^deg. coeff must not be zero
// and deg must not be negative.
func (p *IntPoly) MulMono(q *IntPoly, coeff, deg *big.Int) *IntPoly {
	if coeff.Sign() == 0 {
		panic("zero coefficient")
	}
	if deg.Sign() < 0 {
		panic("negative degree")
	}
	// Since we go left to right, reusing p's term array is okay
	// even if p aliases q.
	terms := p.terms.make(len(q.terms))
	for i, _ := range terms {
		terms[i].Mul(&q.terms[i], coeff, deg)
	}
	p.terms = terms
	return p
}

// Sets p to the product of q and r.
func (p *IntPoly) Mul(q, r *IntPoly) *IntPoly {
	if len(r.terms) > len(q.terms) {
		q, r = r, q
	}
	prod := IntPoly{}
	for _, term := range r.terms {
		t := IntPoly{}
		t.MulMono(q, &term.coeff, &term.deg)
		prod.Add(&prod, &t)
	}
	*p = prod
	return p
}

// Returns a copy of the identity polynomial.
func newIntPolyIdentity() *IntPoly {
	return NewIntPoly([][2]*big.Int{{big.NewInt(1), big.NewInt(0)}})
}

// Sets p to q raised to the kth power. k must be non-negative.
func (p *IntPoly) Pow(q *IntPoly, k *big.Int) *IntPoly {
	if k.Sign() < 0 {
		panic("negative power")
	}
	pow := newIntPolyIdentity()
	for i := k.BitLen() - 1; i >= 0; i-- {
		pow.Mul(pow, pow)
		if k.Bit(i) != 0 {
			pow.Mul(pow, q)
		}
	}
	*p = *pow
	return p
}

// Sets p to q with its coefficients reduced modulo k.
func (p *IntPoly) Mod(q *IntPoly, k *big.Int) *IntPoly {
	if k.Sign() == 0 {
		panic("zero modulus")
	}
	// Since we go left to right, reusing p's term array is okay
	// even if p aliases q.
	terms := p.terms.make(len(q.terms))
	i := 0
	for j, _ := range terms {
		term := intMono{}
		term.Mod(&q.terms[j], k)
		if term.coeff.Sign() != 0 {
			terms[i] = term
			i++
		}
	}
	p.terms = terms[0:i]
	return p
}
