//////////////////////////////////////////////////////////////////////////
// Rust FFI
//
// void execute(
//     unsigned int difficulty,
//     char* input,  int input_size,  /*32*/
//     char* output, int output_size,  /*516...result+proof? wait wut?*/
//     int   sizeInBits
// );

// char /*bool*/ verify(
//     unsigned int  difficulty,
//     char* input,  int input_size,  /*32*/
//     char* output, int output_size, /*516?*/
//     char* proof,  int proof_size,  /*516?*/
//     int   sizeInBits
// );
//////////////////////////////////////////////////////////////////////////

package main

/*
#cgo LDFLAGS: -L./lib -lvdf
#include "./lib/vdf_rust.h"
*/
import "C"

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"unsafe"

	// "vdf_rust"

	"github.com/harmony-one/vdf/src/vdf_go"
)

func main() {
	GenerateVDFAndVerifyGo()
	GenerateVDFAndVerifyRust()
}

const algorithm = "Wesolowski" // or "Pietrzak"
var difficulty = 1000          // low for testing; maybe 10000 production
const inputsize = 32           // match existing go impl
const outputsize = 516         // ^ ... go impl concatenates output and proof here?
const size_in_bits = 2048      // size of long integers in quadratic function group

const bufsize = ((size_in_bits+7)>>3)*2 + 4 // == outputsize

func GenerateVDFAndVerifyRust() {

	input := [inputsize]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	vdf := NewVDFRust(difficulty, input)

	outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	vdf.Execute()

	duration := time.Now().Sub(start)

	output := <-outputChannel

	log.Println(fmt.Sprintf("VDF_Rust computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF_Rust computation finished, time spent %s", duration.String()))

	// var output [outputsize]byte
	// // log.Println(fmt.Sprintf("%s", hex.EncodeToString(output[:])))

	// // ptr := unsafe.Pointer(&input[0])
	// in := (*C.char)(unsafe.Pointer(&input[0]))
	// out := (*C.char)(unsafe.Pointer(&output[0]))

	// // outputChannel := vdf.GetOutputChannel()
	// start := time.Now()

	// C.execute(C.uint(difficulty), in, C.int(inputsize), out, C.int(outputsize), C.int(size_in_bits))

	// duration := time.Now().Sub(start)

	// // output := <-outputChannel
	// //output := input

	// log.Println(fmt.Sprintf("VDF_Rust computation finished, result is  %s", hex.EncodeToString(output[:])))
	// log.Println(fmt.Sprintf("VDF_Rust computation finished, time spent %s", duration.String()))
	// // assert.Equal(t, true, vdf.Verify(output), "failed verifying proof")
}

func GenerateVDFAndVerifyGo() {
	input := [32]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	vdf := vdf_go.New(difficulty, input)

	outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	vdf.Execute()

	duration := time.Now().Sub(start)

	output := <-outputChannel

	log.Println(fmt.Sprintf("VDF_Go computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF_Go computation finished, time spent %s", duration.String()))
}

////////////////////////////////////////////////////
// package vdf_rust {

// VDF is the struct holding necessary state for a hash chain delay function.
type VDFRust struct {
	difficulty int
	input      [inputsize]byte
	output     [516]byte
	outputChan chan [516]byte
	finished   bool
}

// New create a new instance of VDF.
func NewVDFRust(difficulty int, input [inputsize]byte) *VDFRust {
	return &VDFRust{
		difficulty: difficulty,
		input:      input,
		outputChan: make(chan [516]byte),
	}
}

// GetOutputChannel returns the vdf output channel.
// VDF output consists of 258 bytes of serialized Y and  258 bytes of serialized Proof
func (vdf *VDFRust) GetOutputChannel() chan [516]byte {
	return vdf.outputChan
}

// Execute runs the VDF until it's finished and put the result into output channel.
func (vdf *VDFRust) Execute() {
	vdf.finished = false

	in := (*C.char)(unsafe.Pointer(&vdf.input[0]))
	out := (*C.char)(unsafe.Pointer(&vdf.output[0]))

	// outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	C.execute(
		C.uint(difficulty),
		in, C.int(inputsize),
		out, C.int(outputsize),
		C.int(size_in_bits))

	duration := time.Now().Sub(start)

	// output := <-outputChannel

	log.Println(fmt.Sprintf("VDF_Rust computation finished, result is  %s", hex.EncodeToString(vdf.output[:])))
	log.Println(fmt.Sprintf("VDF_Rust computation finished, time spent %s", duration.String()))
	// assert.Equal(t, true, vdf.Verify(output), "failed verifying proof")

	// yBuf, proofBuf := GenerateVDFRust(vdf.input[:], vdf.difficulty, sizeInBits)

	// copy(vdf.output[:], yBuf)
	// copy(vdf.output[258:], proofBuf)

	go func() {
		vdf.outputChan <- vdf.output
	}()

	vdf.finished = true
}

// Verify runs the verification of generated proof
// currently on i7-6700K, verification takes about 350 ms
func (vdf *VDFRust) Verify(proof [258]byte) bool {

	in := (*C.char)(unsafe.Pointer(&vdf.input[0]))
	out := (*C.char)(unsafe.Pointer(&vdf.output[0]))
	prf := (*C.char)(unsafe.Pointer(&proof[0]))

	success := C.verify(
		C.uint(difficulty),
		in, C.int(inputsize),
		out, C.int(outputsize),
		prf, C.int(258),
		C.int(size_in_bits))

	if success == 1 {
		return true
	} else {
		return false
	}
}

// IsFinished returns whether the vdf execution is finished or not.
func (vdf *VDFRust) IsFinished() bool {
	return vdf.finished
}

// GetOutput returns the vdf output, which can be bytes of 0s is the vdf is not finished.
func (vdf *VDFRust) GetOutput() [516]byte {
	return vdf.output
}
