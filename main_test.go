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

func TestRustExampleGo(t *testing.T) {
	// input := [1]byte{0xaa}
	// difficulty := 100
	// expect := "005271e8f9ab2eb8a2906e851dfcb5542e4173f016b85e29d481a108dc82ed3b3f97937b7aa824801138d1771dea8dae2f6397e76a80613afda30f2c30a34b040baaafe76d5707d68689193e5d211833b372a6a4591abb88e2e7f2f5a5ec818b5707b86b8b2c495ca1581c179168509e3593f9a16879620a4dc4e907df452e8dd0ffc4f199825f54ec70472cc061f22eb54c48d6aa5af3ea375a392ac77294e2d955dde1d102ae2ace494293492d31cff21944a8bcb4608993065c9a00292e8d3f4604e7465b4eeefb494f5bea102db343bb61c5a15c7bdf288206885c130fa1f2d86bf5e4634fdc4216bc16ef7dac970b0ee46d69416f9a9acee651d158ac64915b"

	// vdf := vdf_go.New(difficulty, input)

	// outputChannel := vdf.GetOutputChannel()
	// start := time.Now()

	// vdf.Execute()

	// duration := time.Now().Sub(start)

	// var output = <-outputChannel

	// log.Println(fmt.Sprintf("VDF_Go computation finished, result is  %s", hex.EncodeToString(output[:])))
	// log.Println(fmt.Sprintf("VDF_Go computation finished, time spent %s", duration.String()))
	// // assert.Equal(b, true, vdf.Verify(output), "failed verifying proof")
}

func TestGenerateVDFAndVerifyRust(t *testing.T) {

	GenerateVDFAndVerifyRust()

}
func BenchmarkGenerateVDFAndVerifyRust(b *testing.B) {

	GenerateVDFAndVerifyRust()

}
