package main

import (
	"crypto/rand"
	"fmt"
	"tc/config"
	tc "tc/src"
	"time"
)

func main() {
	fmt.Println("Start running timed commitment")

	// group parameter contains prime numbers p,q and a large number n=p*q of groupLength bits
	groupLength := 2048
	groupParameter, err := Setup(groupLength)
	if err != nil {
		fmt.Println("generate group parameter wrong", groupParameter)
	}
	fmt.Println("[Main]N is", groupParameter.N)
	fmt.Println("[Main]primes are", groupParameter.Primes)
	fmt.Println("[Main]", groupParameter.Primes[0].ProbablyPrime(20))
	fmt.Println("[Main]", groupParameter.Primes[1].ProbablyPrime(20))

	// generator g
	timeParameter := 10
	g, mArray, proofSet := tc.GeneratePublicParameter(groupParameter, groupLength, timeParameter)
	fmt.Println("[Main]G is", g)
	fmt.Println("[Main]Length of m array is", len(mArray))

	// setup config
	var mArrayString []string
	for _, m := range mArray {
		mArrayString = append(mArrayString, m.String())
	}

	// write pp to config
	config.SetupConfig(g.String(), groupParameter.N.String(), mArrayString, proofSet)

	// generate commit
	c, h, rKSubOne, rK := tc.GenerateCommit(groupLength, groupParameter)
	fmt.Println("[Main]c is", c)
	fmt.Println("[Main]h is", h)
	fmt.Println("[Main]rKSubOne is", rKSubOne)
	fmt.Println("[Main]rK is", rK)

	// forced open
	start := time.Now()
	tc.ForcedOped(c, h, rKSubOne, rK, timeParameter)
	elapsed := time.Since(start)
	fmt.Println("Time is", elapsed)
}

// Setup Generate Group parameter
func Setup(nLength int) (*tc.GroupParameter, error) {
	groupParameter, err := tc.GenerateGroupParameter(rand.Reader, nLength)

	return groupParameter, err
}
