ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

# LD_LIBRARY_PATH := "$(pwd)/lib"  
# export $LD_LIBRARY_PATH

build:
	cd lib/rust-vdf && cargo build --release
	cp lib/rust-vdf/target/release/librustvdf.so lib/
	go build -ldflags="-r $(ROOT_DIR)lib" main.go
	go build -ldflags="-r $(ROOT_DIR)lib" main_test.go       # just checking

run: build
	./main

test: build	
	go test

bench: build
	go test -bench .