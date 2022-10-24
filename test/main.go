package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Curve is:%v\n", privateKey.PublicKey.Curve.Params())
	msg := "sadfdgrtgtrg"
	hash := sha256.Sum256([]byte(msg))

	for i := 0; i < 20; i++ {
		sk, err := ecdsa.GenerateKey(privateKey.PublicKey.Curve, rand.Reader)
		r, s, err := ecdsa.Sign(rand.Reader, sk, hash[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("signature: %x\n", r)
		fmt.Printf("signature: %x\n", s)
		fmt.Printf("signature: %x\n", *r)
		fmt.Printf("signature: %x\n", *s)

		valid := ecdsa.Verify(&sk.PublicKey, hash[:], r, s)
		fmt.Println("signature verified:", valid)
	}

}
