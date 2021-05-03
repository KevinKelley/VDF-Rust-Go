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

	"github.com/harmony-one/vdf/src/vdf_go"
)

func main() {
	// C.hello(C.CString("John Smith"))

}

func GenerateVDFAndVerifyRust() {
	input := [32]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	// vdf := vdf_go.New(100, input)
	// log.Println(fmt.Sprintf("%s", hex.EncodeToString(input[:])))

	var output [12]byte
	// log.Println(fmt.Sprintf("%s", hex.EncodeToString(output[:])))

	difficulty := 100

	in_size := (C.int)(32)
	out_size := (C.int)(12)

	// ptr := unsafe.Pointer(&input[0])
	in := (*C.char)(unsafe.Pointer(&input[0]))
	out := (*C.char)(unsafe.Pointer(&output[0]))

	// outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	C.execute(C.uint(difficulty), in, in_size, out, out_size, 32)

	duration := time.Now().Sub(start)

	// output := <-outputChannel
	//output := input

	log.Println(fmt.Sprintf("VDF_Go computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF_Go computation finished, time spent %s", duration.String()))
	// assert.Equal(t, true, vdf.Verify(output), "failed verifying proof")
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
	vdf := vdf_go.New(100, input)

	outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	vdf.Execute()

	duration := time.Now().Sub(start)

	output := <-outputChannel

	log.Println(fmt.Sprintf("VDF_Rust computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF_Rust computation finished, time spent %s", duration.String()))
}

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error "+msg+": %v", err)
	}
}

// db, err := sql.Open("mysql", "my_user@/my_database")
// defer db.Close()
// failIf(err, "connecting to my_database")
