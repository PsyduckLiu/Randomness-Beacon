package node

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"entropyNode/config"
	"entropyNode/signature"
	"fmt"
	"math/big"
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
	fmt.Printf("VRF result is:%v\n", hex.EncodeToString(vrfResult))
	fmt.Printf("VRF result is:%v\n", BytesToBinaryString(vrfResult))

	selectBigInt, _ := rand.Int(rand.Reader, big.NewInt(2))
	selectInt, err := strconv.Atoi(selectBigInt.String())
	fmt.Println(selectInt)
	if err == nil {
		panic(err)
	}
	if int(vrfResultBinary[len(vrfResultBinary)-1]) == selectInt {
		fmt.Println("yes")
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
