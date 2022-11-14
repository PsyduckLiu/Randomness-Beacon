package tc

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"tc/config"
)

func GenerateCommit(bits int, groupParameter *GroupParameter) (*big.Int, *big.Int, *big.Int, *big.Int) {
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

	phiN := new(big.Int).Set(bigOne)
	for _, prime := range groupParameter.Primes {
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

	nHash := new(big.Int).SetBytes(Digest(N))
	gHash := new(big.Int).SetBytes(Digest((g)))
	mSubOneHash := new(big.Int).SetBytes(Digest(mArray[len(mArray)-3]))
	mHash := new(big.Int).SetBytes(Digest(mArray[len(mArray)-2]))
	a1Hash := new(big.Int).SetBytes(Digest(a1))
	a2Hash := new(big.Int).SetBytes(Digest(a2))
	a3Hash := new(big.Int).SetBytes(Digest(a3))

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

	result1 := new(big.Int).Set(g)
	result1.Exp(result1, z, N)
	result2 := new(big.Int).Set(h)
	result2.Exp(result2, e, N)
	result1.Mul(result1, result2)
	result1.Mod(result1, N)

	result3 := new(big.Int).Set(mArray[len(mArray)-3])
	result3.Exp(result3, z, N)
	result4 := new(big.Int).Set(rKSubOne)
	result4.Exp(result4, e, N)
	result3.Mul(result3, result4)
	result3.Mod(result3, N)

	result5 := new(big.Int).Set(mArray[len(mArray)-2])
	result5.Exp(result5, z, N)
	result6 := new(big.Int).Set(rK)
	result6.Exp(result6, e, N)
	result5.Mul(result5, result6)
	result5.Mod(result5, N)

	if a1.Cmp(result1) != 0 {
		fmt.Println("test1 error")
	}
	if a2.Cmp(result3) != 0 {
		fmt.Println("test2 error")
	}
	if a3.Cmp(result5) != 0 {
		fmt.Println("test3 error")
	}
	fmt.Println("[Commit]pass all tests!")

	return c, h, rKSubOne, rK
}
