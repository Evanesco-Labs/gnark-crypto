# gurvy

[![License](https://img.shields.io/badge/license-Apache%202-blue)](LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/consensys/gurvy)](https://goreportcard.com/badge/github.com/consensys/gurvy) [![GoDoc](https://godoc.org/github.com/consensys/gurvy?status.svg)](https://godoc.org/github.com/consensys/gurvy)

### Pairing Library implemented in Go ###

`gurvy` implements Elliptic Curve Cryptography (+Pairing) for BLS381, BLS377 and BN256. Originally developed (and used) by [`gnark`](https://github.com/consensys/gnark).

#### Curves supported

* BLS12-381 (Zcash)
* BN256 (Ethereum)
* BLS377 (ZEXE)
* BW6-761 (EC supporting pairing on BLS377 field of definition)

#### Benchmarks

(2,2GHz, i7)

```
BenchmarkPairing-12         1167            954992 ns/op (BLS381)
BenchmarkPairing-12         1534            720039 ns/op (BN256)
BenchmarkPairing-12         1089           1054871 ns/op (BLS377)
BenchmarkPairing-12          355           3341688 ns/op (BW761)
```