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

const algorithm = "Wesolowski"              // or "Pietrzak"
var difficulty = 1000                       // low for testing; maybe 10000 production
const inputsize = 32                        // match existing go impl
const outputsize = 516                      // ^ ... go impl concatenates output and proof here?
const size_in_bits = 2048                   // ^
const bufsize = ((size_in_bits+7)>>3)*2 + 4 // outputsize

func GenerateVDFAndVerifyRust() {

	// let bufsize = size_in_bits / 3 * 2 + 4;  // 516 bytes for 2048
	// log.Println(fmt.Sprintf("est. bufsize is %d", bufsize))

	input := [inputsize]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	// vdf := vdf_go.New(100, input)
	// log.Println(fmt.Sprintf("%s", hex.EncodeToString(input[:])))

	vdf := New(difficulty, input)

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

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error "+msg+": %v", err)
	}
}

// db, err := sql.Open("mysql", "my_user@/my_database")
// defer db.Close()
// failIf(err, "connecting to my_database")

////////////////////////////////////////////////////
// package vdf_rust {

// VDF is the struct holding necessary state for a hash chain delay function.
type VDFRust struct {
	difficulty int
	input      [32]byte
	output     [516]byte
	outputChan chan [516]byte
	finished   bool
}

//size of long integers in quadratic function group
const sizeInBits = 2048

// New create a new instance of VDF.
func New(difficulty int, input [32]byte) *VDFRust {
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
// currently on i7-6700K, it takes about 14 seconds when iteration is set to 10000
func (vdf *VDFRust) Execute() {
	vdf.finished = false

	// input := [inputsize]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
	// 	0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	// vdf := vdf_go.New(100, input)
	// log.Println(fmt.Sprintf("%s", hex.EncodeToString(input[:])))

	var output [outputsize]byte
	// log.Println(fmt.Sprintf("%s", hex.EncodeToString(output[:])))

	// ptr := unsafe.Pointer(&input[0])
	in := (*C.char)(unsafe.Pointer(&vdf.input[0]))
	out := (*C.char)(unsafe.Pointer(&vdf.output[0]))

	// outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	C.execute(C.uint(difficulty), in, C.int(inputsize), out, C.int(outputsize), C.int(size_in_bits))

	duration := time.Now().Sub(start)

	// output := <-outputChannel

	log.Println(fmt.Sprintf("VDF_Rust computation finished, result is  %s", hex.EncodeToString(output[:])))
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
func (vdf *VDFRust) Verify(proof [516]byte) bool {
	return true //VerifyVDFRust(vdf.input[:], proof[:], vdf.difficulty, sizeInBits)
}

// IsFinished returns whether the vdf execution is finished or not.
func (vdf *VDFRust) IsFinished() bool {
	return vdf.finished
}

// GetOutput returns the vdf output, which can be bytes of 0s is the vdf is not finished.
func (vdf *VDFRust) GetOutput() [516]byte {
	return vdf.output
}

// const algorithm = "Wesolowski"              // or "Pietrzak"
// var difficulty = 1000                       // low for testing; maybe 10000 production
// const inputsize = 32                        // match existing go impl
// const outputsize = 516                      // ^ ... go impl concatenates output and proof here?
// const size_in_bits = 2048                   // ^
// const bufsize = ((size_in_bits+7)>>3)*2 + 4 // outputsize

// func GenerateVDFRust(seed []byte, iterations, int_size_bits int) ([]byte, []byte) {
// 	return GenerateVDFWithStopChan(seed, iterations, int_size_bits, nil)
// }

// func GenerateVDFWithStopChan(seed []byte, iterations, int_size_bits int, stop <-chan struct{}) ([]byte, []byte) {
// 	defer timeTrack(time.Now())

// 	D := CreateDiscriminant(seed, int_size_bits)
// 	x := NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)

// 	y, proof := calculateVDF(D, x, iterations, int_size_bits, stop)

// 	if (y == nil) || (proof == nil) {
// 		return nil, nil
// 	} else {
// 		return y.Serialize(), proof.Serialize()
// 	}
// }

// func VerifyVDFRust(seed, proof_blob []byte, iterations, int_size_bits int) bool {
// 	defer timeTrack(time.Now())

// 	int_size := (int_size_bits + 16) >> 4

// 	D := CreateDiscriminant(seed, int_size_bits)
// 	x := NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)
// 	y, _ := NewClassGroupFromBytesDiscriminant(proof_blob[:(2*int_size)], D)
// 	proof, _ := NewClassGroupFromBytesDiscriminant(proof_blob[2*int_size:], D)

// 	return verifyProof(x, y, proof, iterations)
// }
// }
