package main

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	"github.com/pkg/profile"
)

func main() {

	const maxSize = 1 << 22

	pol := make([]fr.Element, maxSize)
	for i := uint64(0); i < maxSize; i++ {
		pol[i].SetRandom()
	}

	domain := fft.NewDomain(uint64(maxSize), 0, false)

	p := profile.Start(profile.TraceProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	domain.FFT(pol, fft.DIF, 0)
	p.Stop()

}
