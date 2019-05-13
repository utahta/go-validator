# go-validator

[![GoDoc Reference](https://godoc.org/github.com/utahta/go-validator?status.svg)](http://godoc.org/github.com/utahta/go-validator)
[![CircleCI](https://circleci.com/gh/utahta/go-validator.svg?style=svg)](https://circleci.com/gh/utahta/go-validator)
[![codecov](https://codecov.io/gh/utahta/go-validator/branch/master/graph/badge.svg)](https://codecov.io/gh/utahta/go-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/utahta/go-validator)](https://goreportcard.com/report/github.com/utahta/go-validator)
[![GitHub release](https://img.shields.io/github/release/utahta/go-validator.svg)](https://github.com/utahta/go-validator/releases)

go-validator is a data validation library for Go.

# Installation

```sh
go get -u github.com/utahta/go-validator
```

# Benchmarks

3.2 GHz Intel Core i7, 64 GB 2667 MHz DDR4
```go
goos: darwin
goarch: amd64
pkg: github.com/utahta/go-validator
BenchmarkValidateVarSuccess-12                          30000000                58.0 ns/op             0 B/op          0 allocs/op
BenchmarkValidateVarParallelSuccess-12                  100000000               12.3 ns/op             0 B/op          0 allocs/op
BenchmarkValidateStructSuccess-12                       10000000               184 ns/op               0 B/op          0 allocs/op
BenchmarkValidateStructParallelSuccess-12               50000000                37.1 ns/op             0 B/op          0 allocs/op
BenchmarkValidateStructComplexSuccess-12                 1000000              1112 ns/op              32 B/op          3 allocs/op
BenchmarkValidateStructComplexParallelSuccess-12        10000000               222 ns/op              32 B/op          3 allocs/op
BenchmarkValidateVarFailure-12                          10000000               187 ns/op             224 B/op          3 allocs/op
BenchmarkValidateStructComplexFailure-12                  500000              3740 ns/op            3760 B/op         60 allocs/op
BenchmarkValidateStructComplexParallelFailure-12         1000000              1117 ns/op            3760 B/op         60 allocs/op
```
