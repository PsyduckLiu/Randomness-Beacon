package tc

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"tc/config"
)

func GenerateCommit(bits int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	mArray := config.GetMArray()
	fmt.Println("[Commit]Length of m array is", len(mArray))
	g := config.GetG()
	fmt.Println("[Commit]Get g is", g)
	N := config.GetN()
	fmt.Println("[Commit]Get N is", N)

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
	fmt.Println("[Commit]h is", h)
	fmt.Println("[Commit]rKSubOne is", rKSubOne)
	fmt.Println("[Commit]rK is", rK)
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
	fmt.Println("[Commit]c is", c)

	return c, h, rKSubOne, rK
}
