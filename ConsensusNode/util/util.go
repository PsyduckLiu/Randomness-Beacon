package util

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// transfer []bytes to string
func BytesToBinaryString(bs []byte) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

// Get Consensus Port(30000 + id)
func PortByID(id int64) int {
	return 30000 + int(id)
}

// Get listening Entropy Port(10000 + id)
func EntropyPortByID(id int64) int {
	return 20000 + int(id)
}

// Hash message v, SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}
