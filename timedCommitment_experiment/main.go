package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"log"
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

	timeParameters := [...]int{15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
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
	c, h, rKSubOne, rK := tc.GenerateCommit(groupLength, groupParameter)
	// fmt.Println("[Main]c is", c)
	// fmt.Println("[Main]h is", h)
	// fmt.Println("[Main]rKSubOne is", rKSubOne)
	// fmt.Println("[Main]rK is", rK)

	// forced open
	start := time.Now()
	tc.ForcedOped(c, h, rKSubOne, rK, timeParameter)
	elapsed := time.Since(start)
	fmt.Println("Time is", elapsed)

	return elapsed
}

// Setup Generate Group parameter
func Setup(nLength int) (*tc.GroupParameter, error) {
	groupParameter, err := tc.GenerateGroupParameter(rand.Reader, nLength)

	return groupParameter, err
}
