# VDF-Rust-Go
Integrate VDF Rust implementation in harmony Go client

## NOTES:

'size_in_bits': requesting 2048.  
In go-vdf, actually get 258 bytes, 2064 bits.  
In rust-vdf, requesting 2048 bits gets result of size 272 bytes or 2176 bits.

In go-vdf, execute() returns a concatenation of 'output' and 'proof', 
258 bytes each, total buf length 516 bytes.

rust-vdf doesn't seem to give access to the 'proof' buffer, 
and returns an output of 136 bytes.  Or 272 bytes.  
I saw both at different times, unclear why


# Rust + Golang
---

This repository shows how, by combining
[`cgo`](https://blog.golang.org/c-go-cgo) and
[Rust's FFI capabilities](https://doc.rust-lang.org/book/ffi.html), we can call
Rust code from Golang.

Run `make build` and then `./main` to see `Rust` + `Golang` in action. You
should see `Hello John Smith!` printed to `stdout`.

## You can do this for your own project
Begin by creating a `lib` directory, where you will keep your Rust libraries.
[Andrew Oppenlander's article on creating a Rust dynamic library is a great introduction](http://oppenlander.me/articles/rust-ffi).

Then, you need to create a C header file for your library. Just copy the `libc`
types that you used.

All that is left to do is to add some `cgo`-specific comments to your Golang
code. These comments tell `cgo` where to find the library and its headers.

```go
/*
#cgo LDFLAGS: -L./lib -lhello
#include "./lib/hello.h"
*/
import "C"
```

> There should not be a newline between `*/` and `import "C"`.

A simple `make build` (use the [Makefile](Makefile) in this repository) will
result in a binary that loads your dynamic library.
