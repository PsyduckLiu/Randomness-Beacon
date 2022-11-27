package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"tc/config"
	tc "tc/src"
	"time"
)

func main() {
	fmt.Println("Start running timed commitment")
	var timeTotal []float64
	var timeAverage []float64

	timeParameters := [...]int{15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}
	for _, timeParameter := range timeParameters {
		time := 0.0

		for i := 0; i < 10; i++ {
			fmt.Println("\n Round", i)
			elapsed := processTC(timeParameter)
			time += elapsed.Seconds()
		}

		timeTotal = append(timeTotal, time)
		timeAverage = append(timeAverage, time/10)
	}

	for index, timeParameter := range timeParameters {
		fmt.Printf("time parameter is %v,total running time is %v\n", timeParameter, timeTotal[index])
		fmt.Printf("time parameter is %v,average running time is %v\n", timeParameter, timeAverage[index])
	}

	writeData(timeAverage, timeParameters[:])
}

func writeData(timeAverage []float64, timeParameters []int) {
	// create a file
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data [][]string
	for index, timeParameter := range timeParameters {
		data = append(data, []string{strconv.FormatFloat(timeAverage[index], 'f', 5, 32), strconv.FormatInt(int64(timeParameter), 10)})
	}

	// write all rows at once
	writer.WriteAll(data)
}

func processTC(timeParameter int) time.Duration {
	// group parameter contains prime numbers p,q and a large number n=p*q of groupLength bits
	groupLength := 2048
	groupParameter, err := Setup(groupLength)
	if err != nil {
		fmt.Println("generate group parameter wrong", groupParameter)
	}
	// fmt.Println("[Main]N is", groupParameter.N)
	// fmt.Println("[Main]primes are", groupParameter.Primes)
	fmt.Println("[Main]Prime?", groupParameter.Primes[0].ProbablyPrime(20))
	fmt.Println("[Main]Prime?", groupParameter.Primes[1].ProbablyPrime(20))

	// generator g
	// timeParameter := 10
	g, mArray, proofSet := tc.GeneratePublicParameter(groupParameter, groupLength, timeParameter)
	// fmt.Println("[Main]G is", g)
	// fmt.Println("[Main]Length of m array is", len(mArray))

	// setup config
	var mArrayString []string
	for _, m := range mArray {
		mArrayString = append(mArrayString, m.String())
	}

	// write pp to config
	config.SetupConfig(g.String(), groupParameter.N.String(), mArrayString, proofSet)

	// generate commit
	c, h, rKSubOne, rK, a1, a2, a3, z := tc.GenerateCommit(groupLength, groupParameter)
	fmt.Println("[Main]c is", c)
	// fmt.Println("[Main]h is", h)
	// fmt.Println("[Main]rKSubOne is", rKSubOne)
	// fmt.Println("[Main]rK is", rK)

	// forced open
	// start := time.Now()
	// tc.ForcedOped(c, h, rKSubOne, rK, timeParameter)
	// elapsed := time.Since(start)
	// fmt.Println("Time is", elapsed)

	// verify tc
	start := time.Now()
	result := VerifyTC(a1.String(), a2.String(), a3.String(), z.String(), h.String(), rKSubOne.String(), rK.String())
	elapsed := time.Since(start)
	if !result {
		fmt.Println("Not pass")
	}
	fmt.Println("Time is", elapsed)

	return elapsed
}

// Setup Generate Group parameter
func Setup(nLength int) (*tc.GroupParameter, error) {
	groupParameter, err := tc.GenerateGroupParameter(rand.Reader, nLength)

	return groupParameter, err
}

// verify TC
func VerifyTC(A1 string, A2 string, A3 string, Z string, H string, RKSubOne string, RK string) bool {
	// fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
	// start := time.Now()
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

	nHash := new(big.Int).SetBytes(Digest(N))
	gHash := new(big.Int).SetBytes(Digest((g)))
	mSubOneHash := new(big.Int).SetBytes(Digest(mArray[len(mArray)-3]))
	mHash := new(big.Int).SetBytes(Digest(mArray[len(mArray)-2]))
	a1Hash := new(big.Int).SetBytes(Digest(a1))
	a2Hash := new(big.Int).SetBytes(Digest(a2))
	a3Hash := new(big.Int).SetBytes(Digest(a3))

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
		return false
	}
	if a2.Cmp(result3) != 0 {
		fmt.Println("===>[VerifyTC]test2 error")
		return false
	}
	if a3.Cmp(result5) != 0 {
		fmt.Println("===>[VerifyTC]test3 error")
		return false
	}

	// end := time.Now()
	// fmt.Println("passed time", end.Sub(start).Seconds())
	return true
}

// Hash message v, using SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}
