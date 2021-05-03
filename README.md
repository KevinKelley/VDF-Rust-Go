# VDF-Rust-Go
Integrate VDF Rust implementation in harmony Go client

## Quickstart
Requires libgmp

- 'make run' or
- 'make test' or
- 'make bench'
- 
## NOTES:

'size_in_bits': requesting 2048.  
In go-vdf, actually get 258 bytes, 2064 bits.  
In rust-vdf (pietrzak algo), requesting 2048 bits gets result of size 272 bytes or 2176 bits.
  - but wesolowski gives 258 bytes for 2048 bitsize
  
In go-vdf, execute() returns a concatenation of 'output' and 'proof', 
258 bytes each, total buf length 516 bytes.

## Parameters
'algorithm': Wesolowski (rust adds Pietrzak)
'difficulty'/'iterations': 100 to test, 10000 maybe
'inputsize': 32 bytes
'outputsize': 516 bytes -- seems to concatenate the output and the proof...
'size_in_bits': 2048 -- could be any power of 2? ...rules?... determines 'outputsize'

## Compatibility
changing any of the Go parameters seems likely to break the chain or at least make cutover tricky...

disappointingly, the timings are ...very sensitive to conditions, I guess.  
At small 'size_in_bits', and small difficulty, I had some runs showing a 500x speedup.
But at the params that the go code uses, that's gone, and the speeds are roughly in the same order of magnitude...
I'm thinking to add the Pietrzak to the benchmarks, in case that one was taking more advantage of the GMP or something...

...now comparison is 2.2s (Go); 827ms (Rust (wesolowski))

