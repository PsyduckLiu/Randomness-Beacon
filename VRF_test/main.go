package main

import (
	"crypto/rand"
	"fmt"
	"vrf/crypto"
)

func main() {
	test()
}

func test() {
	// our "secret keys" are 64 bytes: the spec's 32-byte "secret keys" (which we call the "seed") followed by the 32-byte precomputed public key
	// so the 32-byte "SK" in the test vectors is not directly decoded into a VrfPrivkey, it instead has to go through VrfKeypairFromSeed()

	seed := make([]byte, 32)
	rand.Read(seed)
	fmt.Println("The initial seed is", seed)

	pk, sk := crypto.VrfKeygenFromSeed(seed)
	fmt.Println("Computed public key is", pk)
	fmt.Println("Computed secret key is", sk)

	message := make([]byte, 64)
	fmt.Println("The origin message is", message)

	proof, ok := sk.ProveBytes(message)
	if !ok {
		fmt.Println("Failed to produce a proof")
	} else {
		fmt.Println("proof is", proof)
	}

	ok, output := pk.VerifyBytes(proof, message)
	if !ok {
		fmt.Println("Verify fails on proof", proof)
	} else {
		fmt.Println("vrf output is", output)
	}
}
