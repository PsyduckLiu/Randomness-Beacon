package commitment

import (
	"crypto/rand"
	"entropyNode/config"
	"entropyNode/util"
	"fmt"
	"math/big"
)

var bigOne = big.NewInt(1)

func GenerateTimeCommitment(bits int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {
	mArray := config.GetMArray()
	g := config.GetG()
	N := config.GetN()
	primes := config.GetPrimes()

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

	phiN := new(big.Int).Set(bigOne)
	for _, prime := range primes {
		primeSubOne := new(big.Int)
		primeSubOne.Sub(prime, bigOne)
		fmt.Println(primeSubOne)
		phiN.Mul(phiN, primeSubOne)
	}
	fmt.Println("[Commit]Phi N is", phiN)

	w, _ := rand.Int(rand.Reader, phiN)
	a1 := new(big.Int)
	a1.Exp(g, w, N)
	a2 := new(big.Int)
	a2.Exp(mArray[len(mArray)-3], w, N)
	a3 := new(big.Int)
	a3.Exp(mArray[len(mArray)-2], w, N)

	nHash := new(big.Int).SetBytes(util.Digest(N))
	gHash := new(big.Int).SetBytes(util.Digest((g)))
	mSubOneHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-3]))
	mHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-2]))
	a1Hash := new(big.Int).SetBytes(util.Digest(a1))
	a2Hash := new(big.Int).SetBytes(util.Digest(a2))
	a3Hash := new(big.Int).SetBytes(util.Digest(a3))

	e = big.NewInt(0)
	e.Xor(e, gHash)
	e.Xor(e, nHash)
	e.Xor(e, mSubOneHash)
	e.Xor(e, mHash)
	e.Xor(e, a1Hash)
	e.Xor(e, a2Hash)
	e.Xor(e, a3Hash)

	z := new(big.Int).Set(w)
	alphaE := new(big.Int).Set(e)
	alphaE.Mul(alphaE, alpha)
	z.Sub(z, alphaE)
	z.Mod(z, phiN)

	return c, h, rKSubOne, rK, a1, a2, a3, z
}
