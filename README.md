# gnark-crypto

[![License](https://img.shields.io/badge/license-Apache%202-blue)](LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/consensys/gnark-crypto)](https://goreportcard.com/badge/github.com/consensys/gnark-crypto) [![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/consensys/gnark-crypto)](https://pkg.go.dev/mod/github.com/consensys/gnark-crypto)

`gnark-crypto` provides:
* Elliptic curve cryptography (+pairing) on BN254, BLS381, BLS377 and BW761
* FFT, Polynomial commitment schemes
* MiMC
* EdDSA (on the "companion" twisted edwards curves)

`gnark-crypto` is actively developed and maintained by the team (zkteam@consensys.net) behind:
* [`gnark`: a framework to execute (and verify) algorithms in zero-knowledge](https://github.com/consensys/gnark) 


## Warning
**`gnark-crypto` has not been audited and is provided as-is, use at your own risk. In particular, `gnark-crypto` makes no security guarantees such as constant time implementation or side-channel attack resistance.**

`gnark-crypto` packages are optimized for 64bits architectures (x86 `amd64`) and tested on Unix (Linux / macOS).


## Getting started

### Go version

`gnark-crypto` is tested with the last 2 major releases of Go (1.15 and 1.16).

### Install `gnark-crypto` 

```bash
go get github.com/consensys/gnark-crypto
```

Note if that if you use go modules, in `go.mod` the module path is case sensitive (use `consensys` and not `ConsenSys`).

### Documentation

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/consensys/gnark-crypto)](https://pkg.go.dev/mod/github.com/consensys/gnark-crypto)

The APIs are consistent accross the curves. For example, [here is `bn254` godoc](https://pkg.go.dev/github.com/consensys/gnark-crypto/ecc/bn254#pkg-overview).

## Benchmarks

[Benchmarking pairing-friendly elliptic curves libraries](https://hackmd.io/@zkteam/eccbench) 

>The libraries are implemented in different languages and some use more assembly code than others. Besides the different algorithmic and software optimizations used across, it should be noted also that some libraries target constant-time implementation for some operations making it de facto slower. However, it can be clear that consensys/gnark-crypto is one of the fastest pairing-friendly elliptic curve libraries to be used in zkp projects with different curves.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/consensys/gnark-crypto/tags). 


## License

This project is licensed under the Apache 2 License - see the [LICENSE](LICENSE) file for details
