package util

import (
	"bytes"
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

// Get Port(10000 + id)
func EntropyPortByID(id int) int {
	return 20000 + int(id)
}
