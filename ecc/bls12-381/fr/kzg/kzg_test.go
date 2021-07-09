// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package kzg

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
)

// newScheme returns a new KZG scheme.
// This should be used for testing purpose only
// it creates a new FFT domain and a new SRS without randomness
func newScheme(size int) *Scheme {
	return &Scheme{
		SRS:    NewSRS(size, new(big.Int).SetInt64(42)),
		Domain: fft.NewDomain(uint64(size), 0, false),
	}
}

func randomPolynomial(size int) polynomial.Polynomial {
	f := make(polynomial.Polynomial, size)
	for i := 0; i < size; i++ {
		f[i].SetRandom()
	}
	return f
}

func TestDividePolyByXminusA(t *testing.T) {

	sizePol := 230

	domain := fft.NewDomain(uint64(sizePol), 0, false)

	// build random polynomial
	pol := make(polynomial.Polynomial, sizePol)
	for i := 0; i < sizePol; i++ {
		pol[i].SetRandom()
	}

	// evaluate the polynomial at a random point
	var point fr.Element
	point.SetRandom()
	evaluation := pol.Eval(&point)

	// compute f-f(a)/x-a
	h := dividePolyByXminusA(domain, pol, evaluation, point)

	if len(h) != 229 {
		t.Fatal("inconsistant size of quotient")
	}

	// probabilistic test (using Schwartz Zippel lemma, evaluation at one point is enough)
	var randPoint, xminusa fr.Element
	randPoint.SetRandom()

	polRandpoint := pol.Eval(&randPoint)
	polRandpoint.Sub(&polRandpoint, &evaluation) // f(rand)-f(point)

	hRandPoint := h.Eval(&randPoint)
	xminusa.Sub(&randPoint, &point) // rand-point

	// f(rand)-f(point)	==? h(rand)*(rand-point)
	hRandPoint.Mul(&hRandPoint, &xminusa)

	if !hRandPoint.Equal(&polRandpoint) {
		t.Fatal("Error f-f(a)/x-a")
	}
}

func TestSerialization(t *testing.T) {

	// create a KZG scheme
	srs := NewSRS(64, new(big.Int).SetInt64(42))

	// serialize it...
	var buf bytes.Buffer
	_, err := srs.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	// reconstruct the scheme
	var _srs SRS
	_, err = _srs.ReadFrom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	// compare
	if !reflect.DeepEqual(srs, &_srs) {
		t.Fatal("scheme serialization failed")
	}

}

func TestCommit(t *testing.T) {

	// create a KZG scheme
	s := newScheme(64)

	// create a polynomial
	f := make(polynomial.Polynomial, 60)
	for i := 0; i < 60; i++ {
		f[i].SetRandom()
	}

	// commit using the method from KZG
	_kzgCommit, err := s.Commit(f)
	if err != nil {
		t.Fatal(err)
	}
	var kzgCommit bls12381.G1Affine
	kzgCommit.Unmarshal(_kzgCommit.Marshal())

	// check commitment using manual commit
	var x fr.Element
	x.SetString("42")
	fx := f.Eval(&x)
	var fxbi big.Int
	fx.ToBigIntRegular(&fxbi)
	var manualCommit bls12381.G1Affine
	manualCommit.Set(&s.SRS.G1[0])
	manualCommit.ScalarMultiplication(&manualCommit, &fxbi)

	// compare both results
	if !kzgCommit.Equal(&manualCommit) {
		t.Fatal("error KZG commitment")
	}

}

func TestVerifySinglePoint(t *testing.T) {

	// create a KZG scheme
	s := newScheme(64)

	// create a polynomial
	f := randomPolynomial(60)

	// commit the polynomial
	digest, err := s.Commit(f)
	if err != nil {
		t.Fatal(err)
	}

	// compute opening proof at a random point
	var point fr.Element
	point.SetString("4321")
	proof, err := s.Open(&point, f)
	if err != nil {
		t.Fatal(err)
	}

	// verify the claimed valued
	expected := f.Eval(&point)
	if !proof.ClaimedValue.Equal(&expected) {
		t.Fatal("inconsistant claimed value")
	}

	// verify correct proof
	err = s.Verify(&digest, &proof)
	if err != nil {
		t.Fatal(err)
	}

	// verify wrong proof
	proof.ClaimedValue.Double(&proof.ClaimedValue)
	err = s.Verify(&digest, &proof)
	if err == nil {
		t.Fatal("verifying wrong proof should have failed")
	}
}

func TestBatchVerifySinglePoint(t *testing.T) {

	// create a KZG scheme
	s := newScheme(64)

	// create polynomials
	f := make([]polynomial.Polynomial, 10)
	for i := 0; i < 10; i++ {
		f[i] = randomPolynomial(60)
	}

	// commit the polynomials
	digests := make([]Digest, 10)
	for i := 0; i < 10; i++ {
		digests[i], _ = s.Commit(f[i])

	}

	// compute opening proof at a random point
	var point fr.Element
	point.SetString("4321")
	proof, err := s.BatchOpenSinglePoint(&point, digests, f)
	if err != nil {
		t.Fatal(err)
	}

	// verify the claimed values
	for i := 0; i < 10; i++ {
		expectedClaim := f[i].Eval(&point)
		if !expectedClaim.Equal(&proof.ClaimedValues[i]) {
			t.Fatal("inconsistant claimed values")
		}
	}

	// verify correct proof
	err = s.BatchVerifySinglePoint(digests, &proof)
	if err != nil {
		t.Fatal(err)
	}

	// verify wrong proof
	proof.ClaimedValues[0].Double(&proof.ClaimedValues[0])
	err = s.BatchVerifySinglePoint(digests, &proof)
	if err == nil {
		t.Fatal("verifying wrong proof should have failed")
	}

}

func TestBatchVerifyMultiPoints(t *testing.T) {

	// create a KZG scheme
	s := newScheme(64)

	// create polynomials
	f := make([]polynomial.Polynomial, 10)
	for i := 0; i < 10; i++ {
		f[i] = randomPolynomial(60)
	}

	// commit the polynomials
	digests := make([]Digest, 10)
	for i := 0; i < 10; i++ {
		digests[i], _ = s.Commit(f[i])

	}

	// compute 2 batch opening proofs at 2 random points
	points := make([]fr.Element, 2)
	batchProofs := make([]BatchOpeningProof, 2)
	points[0].SetRandom()
	batchProofs[0], _ = s.BatchOpenSinglePoint(&points[0], digests[:5], f[:5])
	points[1].SetRandom()
	batchProofs[1], _ = s.BatchOpenSinglePoint(&points[1], digests[5:], f[5:])

	// fold the 2 batch opening proofs
	proofs := make([]OpeningProof, 2)
	foldedDigests := make([]Digest, 2)
	proofs[0], foldedDigests[0], _ = s.FoldProof(digests[:5], &batchProofs[0])
	proofs[1], foldedDigests[1], _ = s.FoldProof(digests[5:], &batchProofs[1])

	// check the the individual batch proofs are correct
	err := s.Verify(&foldedDigests[0], &proofs[0])
	if err != nil {
		t.Fatal(err)
	}
	err = s.Verify(&foldedDigests[1], &proofs[1])
	if err != nil {
		t.Fatal(err)
	}

	// batch verify correct folded proofs
	err = s.BatchVerifyMultiPoints(foldedDigests, proofs)
	if err != nil {
		t.Fatal(err)
	}

	// batch verify tampered folded proofs
	proofs[0].ClaimedValue.Double(&proofs[0].ClaimedValue)
	err = s.BatchVerifyMultiPoints(foldedDigests, proofs)
	if err == nil {
		t.Fatal(err)
	}

}

const benchSize = 1 << 16

func BenchmarkKZGCommit(b *testing.B) {
	// kzg scheme
	s := newScheme(benchSize)

	// random polynomial
	p := randomPolynomial(benchSize / 2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Commit(p)
	}
}

func BenchmarkKZGOpen(b *testing.B) {
	// kzg scheme
	s := newScheme(benchSize)

	// random polynomial
	p := randomPolynomial(benchSize / 2)
	var r fr.Element
	r.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Open(&r, p)
	}
}

func BenchmarkKZGVerify(b *testing.B) {
	// kzg scheme
	s := newScheme(benchSize)

	// random polynomial
	p := randomPolynomial(benchSize / 2)
	var r fr.Element
	r.SetRandom()

	// commit
	comm, err := s.Commit(p)
	if err != nil {
		b.Fatal(err)
	}

	// open
	openingProof, err := s.Open(&r, p)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Verify(&comm, &openingProof)
	}
}

func BenchmarkKZGBatchOpen10(b *testing.B) {
	// kzg scheme
	s := newScheme(benchSize)

	// 10 random polynomials
	var ps [10]polynomial.Polynomial
	for i := 0; i < 10; i++ {
		ps[i] = randomPolynomial(benchSize / 2)
	}

	// commitments
	var commitments [10]Digest
	for i := 0; i < 10; i++ {
		commitments[i], _ = s.Commit(ps[i])
	}

	var r fr.Element
	r.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.BatchOpenSinglePoint(&r, commitments[:], ps[:])
	}
}

func BenchmarkKZGBatchVerify10(b *testing.B) {
	// kzg scheme
	s := newScheme(benchSize)

	// 10 random polynomials
	var ps [10]polynomial.Polynomial
	for i := 0; i < 10; i++ {
		ps[i] = randomPolynomial(benchSize / 2)
	}

	// commitments
	var commitments [10]Digest
	for i := 0; i < 10; i++ {
		commitments[i], _ = s.Commit(ps[i])
	}

	var r fr.Element
	r.SetRandom()

	proof, err := s.BatchOpenSinglePoint(&r, commitments[:], ps[:])
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.BatchVerifySinglePoint(commitments[:], &proof)
	}
}
