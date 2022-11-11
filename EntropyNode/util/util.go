package util

import (
	"bytes"
	"crypto/sha256"

	"github.com/algorand/go-algorand/crypto"

	"fmt"
)

// transfer []bytes to string
// func BytesToBinaryString(bs []byte) string {
// 	buf := bytes.NewBuffer([]byte{})
// 	for _, v := range bs {
// 		buf.WriteString(fmt.Sprintf("%08b", v))
// 	}
// 	return buf.String()
// }

func BytesToBinaryString(bs crypto.VrfProof) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

// Get Port(10000 + id)
func EntropyPortByID(id int) int {
	return 20000 + int(id)
}

// Hash message v, SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}
