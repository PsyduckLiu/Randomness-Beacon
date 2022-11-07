package commitment

import (
	"crypto/rand"
	"entropyNode/config"
	"fmt"
	"math/big"
)

// func GenerateTimeCommitment() *big.Int {
// 	var upper, e = big.NewInt(2), big.NewInt(256)
// 	upper.Exp(upper, e, nil)

// 	// generate random number
// 	tc, err := rand.Int(rand.Reader, upper)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return tc
// }

func GenerateTimeCommitment(bits int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	mArray := config.GetMArray()
	g := config.GetG()
	N := config.GetN()

	nSqrt := new(big.Int)
	upperBound := new(big.Int)
	nSqrt.Sqrt(N)
	nSqrt.Add(nSqrt, nSqrt)
	upperBound.Sub(N, nSqrt)
	fmt.Println("[Commit]upper bound is", upperBound)

	alpha, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		fmt.Println("[Commit]generate alpha wrong", err)
	}
	fmt.Println("[Commit]alpha is", alpha)

	h := new(big.Int)
	rKSubOne := new(big.Int)
	rK := new(big.Int)
	r := new(big.Int)
	h.Exp(g, alpha, N)
	rKSubOne.Exp(mArray[len(mArray)-3], alpha, N)
	rK.Exp(mArray[len(mArray)-2], alpha, N)
	r.Exp(mArray[len(mArray)-1], alpha, N)
	fmt.Println("[Commit]r is", r)

	c := new(big.Int)
	upper, e := big.NewInt(2), big.NewInt(int64(bits))
	upper.Exp(upper, e, nil)

	// generate random number
	msg, err := rand.Int(rand.Reader, upper)
	if err != nil {
		panic(err)
	}

	c.Xor(msg, r)
	fmt.Println("[Commit]msg is", msg)

	return c, h, rKSubOne, rK
}
