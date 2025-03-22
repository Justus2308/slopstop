package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestShortInput(t *testing.T) {
	input := "flying fish flew by the space station"
	expectedShingles := map[string]uint{
		"ta": 0,
		"e ": 0,
		"y ": 0,
		"fi": 0,
		"ti": 0,
		"io": 0,
		" s": 0,
		"is": 0,
		"fl": 0,
		"sp": 0,
		" t": 0,
		"yi": 0,
		"by": 0,
		"ce": 0,
		"th": 0,
		"at": 0,
		" f": 0,
		"ly": 0,
		" b": 0,
		"g ": 0,
		"st": 0,
		"w ": 0,
		"sh": 0,
		"pa": 0,
		"ng": 0,
		"ac": 0,
		"h ": 0,
		"le": 0,
		"in": 0,
		"on": 0,
		"he": 0,
		"ew": 0,
	}
	// input := "test"
	// expected_shingles := set[string]{
	// 	"te": 0,
	// 	"es": 0,
	// 	"st": 0,
	// }
	reader := strings.NewReader(input)
	r := bufio.NewReader(reader)
	hasher, _ := NewHasher(2)
	shinglesAdded, err := hasher.AddShingles(r)
	if err != nil {
		t.Fatalf("AddShingles failed: %s", err.Error())
	}
	shingles := hasher.vocab
	for k := range shingles {
		t.Logf("%s\n", k)
	}
	is_eql_shingles := (shinglesAdded == len(shingles) && len(expectedShingles) == len(shingles))
	for k := range expectedShingles {
		_, ok := shingles[k]
		if !ok {
			t.Logf("missing shingle: %s", k)
		}
		is_eql_shingles = is_eql_shingles && ok
	}
	if !is_eql_shingles {
		t.Fatal("AddShingles produced incorrect shingles")
	}

	input1 := "he will not allow you to bring your sticks of dynamite and pet armadillo along"
	input2 := "he figured a few sticks of dynamite were easier than a fishing pole to catch an armadillo"
	reader1 := strings.NewReader(input1)
	reader2 := strings.NewReader(input2)
	r1 := bufio.NewReader(reader1)
	r2 := bufio.NewReader(reader2)
	_, err = hasher.AddShingles(r1)
	if err != nil {
		t.Fatalf("AddShingles failed: %s", err.Error())
	}
	_, err = hasher.AddShingles(r2)
	if err != nil {
		t.Fatalf("AddShingles failed: %s", err.Error())
	}

	hasher.Finalize(20)

	signature := hasher.Hash(input)
	t.Logf("signature: %v", signature)

	signature1 := hasher.Hash(input1)
	t.Logf("signature1: %v", signature1)

	signature2 := hasher.Hash(input2)
	t.Logf("signature2: %v", signature2)

	sim01 := Similiarity(signature, signature1)
	t.Logf("sim01: %f", sim01)

	sim02 := Similiarity(signature, signature2)
	t.Logf("sim02: %f", sim02)

	sim12 := Similiarity(signature1, signature2)
	t.Logf("sim12: %f", sim12)

	if sim12 < sim01 || sim12 < sim02 {
		t.Fatalf("Signatures produced incorrect similiarity")
	}
}
