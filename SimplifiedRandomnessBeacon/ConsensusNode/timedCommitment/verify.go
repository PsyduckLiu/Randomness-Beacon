package tc

import (
	"consensusNode/config"
	"consensusNode/util"
	"fmt"
	"math/big"
	"runtime"
	"time"
)

// verify TC
func VerifyTC(A1 string, A2 string, A3 string, Z string, H string, RKSubOne string, RK string) (bool, float64) {
	fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
	start := time.Now()
	mArray := config.GetMArray()
	g := config.GetG()
	N := config.GetN()

	a1, _ := new(big.Int).SetString(A1, 10)
	a2, _ := new(big.Int).SetString(A2, 10)
	a3, _ := new(big.Int).SetString(A3, 10)
	z, _ := new(big.Int).SetString(Z, 10)
	h, _ := new(big.Int).SetString(H, 10)
	rKSubOne, _ := new(big.Int).SetString(RKSubOne, 10)
	rK, _ := new(big.Int).SetString(RK, 10)

	nHash := new(big.Int).SetBytes(util.Digest(N))
	gHash := new(big.Int).SetBytes(util.Digest((g)))
	mSubOneHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-3]))
	mHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-2]))
	a1Hash := new(big.Int).SetBytes(util.Digest(a1))
	a2Hash := new(big.Int).SetBytes(util.Digest(a2))
	a3Hash := new(big.Int).SetBytes(util.Digest(a3))

	e := big.NewInt(0)
	e.Xor(e, gHash)
	e.Xor(e, nHash)
	e.Xor(e, mSubOneHash)
	e.Xor(e, mHash)
	e.Xor(e, a1Hash)
	e.Xor(e, a2Hash)
	e.Xor(e, a3Hash)
	fmt.Println("after xor", e)

	// start1 := time.Now()
	result1 := new(big.Int).Set(g)
	result1.Exp(result1, z, N)
	result2 := new(big.Int).Set(h)
	result2.Exp(result2, e, N)
	result1.Mul(result1, result2)
	result1.Mod(result1, N)
	// end1 := time.Now()
	// fmt.Println("passed time1", end1.Sub(start1).Seconds())

	// start2 := time.Now()
	result3 := new(big.Int).Set(mArray[len(mArray)-3])
	result3.Exp(result3, z, N)
	result4 := new(big.Int).Set(rKSubOne)
	result4.Exp(result4, e, N)
	result3.Mul(result3, result4)
	result3.Mod(result3, N)
	// end2 := time.Now()
	// fmt.Println("passed time2", end2.Sub(start2).Seconds())

	// start3 := time.Now()
	result5 := new(big.Int).Set(mArray[len(mArray)-2])
	result5.Exp(result5, z, N)
	result6 := new(big.Int).Set(rK)
	result6.Exp(result6, e, N)
	result5.Mul(result5, result6)
	result5.Mod(result5, N)
	// end3 := time.Now()
	// fmt.Println("passed time3", end3.Sub(start3).Seconds())

	if a1.Cmp(result1) != 0 {
		fmt.Println("===>[VerifyTC]test1 error")
		return false, 0.0
	}
	if a2.Cmp(result3) != 0 {
		fmt.Println("===>[VerifyTC]test2 error")
		return false, 0.0
	}
	if a3.Cmp(result5) != 0 {
		fmt.Println("===>[VerifyTC]test3 error")
		return false, 0.0
	}

	end := time.Now()
	fmt.Println("passed time", end.Sub(start).Seconds())
	return true, end.Sub(start).Seconds()
}
