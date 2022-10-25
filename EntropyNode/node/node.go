package node

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"entropyNode/config"
	"entropyNode/signature"
	"fmt"
	"strconv"
)

// initialize an entropy node
func StartEntropyNode(id int) {
	fmt.Printf("[Node%d] is running\n", id)

	config.ReadConfig()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey(marshalledCurve)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Key message parse err:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve
	fmt.Printf("Curve is %v\n", curve.Params())

	// generate private key
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR]Generate private key err:%s", err))
	}
	fmt.Printf("===>My own key is: %v\n", privateKey)

	// calculate VRF result
	previousOutput := config.GetPreviousInput()
	vrfResult := calVRF(previousOutput, privateKey)
	vrfResultBinary := BytesToBinaryString(vrfResult)
	fmt.Printf("VRF result is:%v\n", BytesToBinaryString(vrfResult))
	fmt.Printf("VRF result last bit is:%v\n", vrfResultBinary[len(vrfResultBinary)-1:])

	// match VRF result with difficulty
	difficulty := config.GetDifficulty()
	vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
	if err != nil {
		panic(err)
	}
	if vrfResultTail == difficulty {
		fmt.Println("yes!!!!!!!!!!")
	}
}

// calculate VRF output
func calVRF(previousOutput []byte, sk *ecdsa.PrivateKey) []byte {
	// skD := sk.D
	vrfRes := signature.GenerateSig(previousOutput, sk)

	valid := signature.VerifySig(previousOutput, vrfRes, &sk.PublicKey)
	fmt.Printf("Verify result is:%v\n", valid)
	return vrfRes
}

func BytesToBinaryString(bs []byte) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}
