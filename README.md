# GoFprint: A Go library to work with [libfprint](https://github.com/freedesktop/libfprint) via [cgo](https://github.com/golang/go/wiki/cgo)

![gopher tatooed](https://image.ibb.co/feRi0q/rsz-1castor.png)

This package comes with functions that wrap libfprint's C functions in order to work directly with fingerprint devices (tested only with URU4000)
while using Go's garbage collection features and its data structures like []byte.

## Installation

To install this package, all you have to do is run the following command.

```
$ go get github.com/caws/gofprint
```

## Get started

**Before anything, remember to install libfprint following the instructions otherwise you won't be able to use the bindings.**

You'll find more information in the examples folder, however, I'll be updating this section with explanations about how to use these bindings sometime in the future.

## Examples

In the [examples](examples) folder you can find two examples for single fingerprint verification and multiple fingerprint verification.

To run the example, do the following:

```
$ cd gofprint/examples/singlePrint
$ go build singlePrint.go
$ sudo ./singleFprint
```

PS: You need to use the sudo command to allow access the device.

## Issues

I tried to clean up the code and get rid of all bugs I could before releasing it as a package,
however, bugs are stubborn and sometimes hide in plain sight.

So, in case you find a bug, let me know or create a pull request. I'll be more than happy to review and accept it if need be.

Charles Washington