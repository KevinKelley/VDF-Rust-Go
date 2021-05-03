package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/harmony-one/vdf/src/vdf_go"
)

func BenchmarkGenerateVDFAndVerifyGo(b *testing.B) {
	input := [32]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	vdf := vdf_go.New(100, input)

	outputChannel := vdf.GetOutputChannel()
	start := time.Now()

	vdf.Execute()

	duration := time.Now().Sub(start)

	var output = <-outputChannel

	log.Println(fmt.Sprintf("VDF_Go computation finished, result is  %s", hex.EncodeToString(output[:])))
	log.Println(fmt.Sprintf("VDF_Go computation finished, time spent %s", duration.String()))
	// assert.Equal(b, true, vdf.Verify(output), "failed verifying proof")
}

func TestGenerateVDFAndVerifyRust(t *testing.T) {

	GenerateVDFAndVerifyRust()

}
func BenchmarkGenerateVDFAndVerifyRust(b *testing.B) {

	GenerateVDFAndVerifyRust()

}
