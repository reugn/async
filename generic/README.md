# Go 1.18 Generics package
This package contains preliminary [`Go Generics`](https://github.com/golang/go/issues/43651) prototypes.

## Implemented data types
* asyncTask <sup>[1](#unexported)</sup>

<sup name="unexported">1</sup>  Unexported since it is not possible to export generic code yet.

## Getting started

### Go 1.17
```sh
go test ./... -gcflags=-G=3 -vet=off
```

### Go 1.18+
Install [gotip](https://pkg.go.dev/golang.org/dl/gotip):
```sh
go install golang.org/dl/gotip@latest
gotip download
```
Use `gotip` to build and test the code.
