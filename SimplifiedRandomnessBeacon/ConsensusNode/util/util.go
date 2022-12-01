package util

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/algorand/go-algorand/crypto"
)

// to be modified
const MaxFaultyNode = 2
const TotalNodeNum = 3*MaxFaultyNode + 1

// convert crypto.VrfProof([80]byte) to binary string
func BytesToBinaryString(bs crypto.VrfProof) string {
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

// Get listening Entropy Port(20000 + id)
func EntropyPortByID(id int64) int {
	return 20000 + int(id)
}

// Hash message v, using SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}

// wriete output to output.txt
func WriteResult(output string) {
	filePath := "../output.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("file failed", err)
		}

	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(time.Now().String() + "\n")
	write.WriteString(output + "\n")

	write.Flush()
}
