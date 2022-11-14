package main

import (
	"crypto/rand"
	"fmt"

	"github.com/gatechain/crypto"
)

type TestingHashable struct {
	data []byte
}

func (s TestingHashable) ToBeHashed() (crypto.HashID, []byte) {
	return "test", s.data
}

func randString() (b TestingHashable) {
	d := make([]byte, 100)
	_, err := rand.Read(d)
	if err != nil {
		panic(err)
	}
	return TestingHashable{d}
}

func main() {
	msg := randString()

	pk, sk := crypto.VrfKeygen()
	fmt.Println("public key is", pk)
	fmt.Println("secret key is", sk)

	proof, ok := sk.Prove(msg)
	if !ok {
		panic("Failed to construct VRF proof")
	}
	fmt.Println("VRF Proof is", proof)

	ok, output := pk.Verify(proof, msg)
	if !ok {
		fmt.Println("Verify() fails on proof")
	}
	fmt.Println("VRF Output is", output)
}
