package main

import (
	"fmt"

	"github.com/algorand/go-algorand/crypto"
)

func main() {
	pk, sk := crypto.VrfKeygen()
	fmt.Println("public key is", pk)
	fmt.Println("secret key is", sk)

	// vrfPrivKey := toVrfPrivKey(os.Getenv("VRFPRIV"))

	// msg := getRoundSeedHashable(os.Args[2], os.Args[3])

	// vrfProof, ok := vrfPrivKey.Prove(msg)

	// if !ok {
	// 	panic("Proof failed.")
	// }

	// ok1, output := getPublicKey(os.Args[1]).Verify(vrfProof, Msg(msg))
	// if !ok1 {
	// 	panic("Verification failed.")
	// }

	// fmt.Println(base32.StdEncoding.EncodeToString(vrfProof[:]),
	// 	base32.StdEncoding.EncodeToString(output[:]))
}
