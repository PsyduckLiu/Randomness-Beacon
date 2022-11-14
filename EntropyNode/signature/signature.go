package signature

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

// hash the input message, using sSHA256
func Digest(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)
	digest := h.Sum(nil)

	return digest
}

// generate the signature of msg
func GenerateSig(msg []byte, sk *ecdsa.PrivateKey) []byte {
	digest := Digest(msg)
	sig, err := ecdsa.SignASN1(rand.Reader, sk, digest)
	if err != nil {
		panic(err)
	}

	return sig
}

// verify the signature of msg
func VerifySig(msg []byte, sig []byte, pk *ecdsa.PublicKey) bool {
	digest := Digest(msg)
	valid := ecdsa.VerifyASN1(pk, digest, sig)

	return valid
}
