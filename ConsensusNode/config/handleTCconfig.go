package config

import (
	"consensusNode/util"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/spf13/viper"
)

var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

type GroupParameter struct {
	N      *big.Int
	Primes []*big.Int
}

// generate group parameters for timed commitment
func GenerateGroupParameter(random io.Reader, bits int) (*GroupParameter, error) {
	MaybeReadByte(random)
	groupP := new(GroupParameter)

	if bits < 64 {
		primeLimit := float64(uint64(1) << uint(bits/2))
		// pi approximates the number of primes less than primeLimit
		pi := primeLimit / (math.Log(primeLimit) - 1)
		// Generated primes start with 11 (in binary) so we can only
		// use a quarter of them.
		pi /= 4
		// Use a factor of two to ensure that key generation terminates
		// in a reasonable amount of time.
		pi /= 2
		if pi <= float64(2) {
			return nil, errors.New("crypto/group: too few primes of given length to generate an group key")
		}
	}

	primes := make([]*big.Int, 2)

NextSetOfPrimes:
	for {
		todo := bits
		// crypto/rand should set the top two bits in each prime.
		// Thus each prime has the form
		//   p_i = 2^bitten(p_i) × 0.11... (in base 2).
		// And the product is:
		//   P = 2^todo × α
		// where α is the product of 2 numbers of the form 0.11...
		//
		// If α < 1/2 (which can happen for 2 > 2), we need to
		// shift todo to compensate for lost bits: the mean value of 0.11...
		// is 7/8, so todo + shift - 2 * log2(7/8) ~= bits - 1/2
		// will give good results.
		for i := 0; i < 2; i++ {
			var err error
			primes[i], err = rand.Prime(random, todo/(2-i))
			if err != nil {
				return nil, err
			}
			todo -= primes[i].BitLen()
		}

		// Make sure that primes is pairwise unequal.
		if primes[0].Cmp(primes[1]) == 0 {
			continue NextSetOfPrimes
		}

		n := new(big.Int).Set(bigOne)
		for _, prime := range primes {
			n.Mul(n, prime)
		}
		if n.BitLen() != bits {
			// This should never happen for 2 == 2 because
			// crypto/rand should set the top two bits in each prime.
			// For 2 > 2 we hope it does not happen often.
			continue NextSetOfPrimes
		}

		groupP.Primes = primes
		groupP.N = n
		break
	}

	return groupP, nil
}

// generate public parameters for timed commitment
func GeneratePublicParameter(groupP *GroupParameter, bits int, k int) (*big.Int, []*big.Int, [][3]string) {
	var mArray []*big.Int
	primesUnder128 := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127}

	// generate h
	h, err := rand.Int(rand.Reader, groupP.N)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GeneratePublicParameter]Generate h failed:%s", err))
	}
	fmt.Println("===>[Setup]H is", h)

	// calculate phiN
	phiN := new(big.Int).Set(bigOne)
	for _, prime := range groupP.Primes {
		primeSubOne := new(big.Int)
		primeSubOne.Sub(prime, bigOne)
		phiN.Mul(phiN, primeSubOne)
	}
	fmt.Println("===>[Setup]PhiN is", phiN)

	// calculate g
	power := new(big.Int).Set(bigOne)
	for _, p := range primesUnder128 {
		q := new(big.Int)
		q.Exp(big.NewInt(p), big.NewInt(int64(bits)), phiN)
		power.Mul(power, q)
		power.Mod(power, phiN)
	}
	g := new(big.Int)
	g.Exp(h, power, groupP.N)
	fmt.Println("===>[Setup]G is", g)

	// calculate mArray
	for i := 0; i <= k; i++ {
		power1 := new(big.Int)
		power2 := new(big.Int)
		m := new(big.Int)
		power1.Exp(bigTwo, big.NewInt(int64(i)), nil)
		power2.Exp(bigTwo, power1, phiN)
		m.Exp(g, power2, groupP.N)
		mArray = append(mArray, m)

		if i == k {
			mFirst := new(big.Int)
			power1.Sub(power1, bigOne)
			power2.Exp(bigTwo, power1, phiN)
			mFirst.Exp(g, power2, groupP.N)
			mArray = append(mArray, mFirst)
		}
	}
	fmt.Println("===>[Setup]Length of m array is", len(mArray))

	// calculate proofSet
	var proofSet [][3]string
	for i := 0; i < k; i++ {
		alpha, _ := rand.Int(rand.Reader, phiN)
		z := new(big.Int)
		z.Exp(g, alpha, groupP.N)
		w := new(big.Int)
		w.Exp(mArray[i], alpha, groupP.N)

		gHash := new(big.Int).SetBytes(util.Digest((g)))
		nHash := new(big.Int).SetBytes(util.Digest(groupP.N))
		biHash := new(big.Int).SetBytes(util.Digest(mArray[i]))
		biPlusHash := new(big.Int).SetBytes(util.Digest(mArray[i+1]))
		zHash := new(big.Int).SetBytes(util.Digest(z))
		wHash := new(big.Int).SetBytes(util.Digest(w))

		c := big.NewInt(0)
		c.Xor(c, gHash)
		c.Xor(c, nHash)
		c.Xor(c, biHash)
		c.Xor(c, biPlusHash)
		c.Xor(c, zHash)
		c.Xor(c, wHash)

		y := new(big.Int).Set(c)
		power3 := new(big.Int)
		power4 := new(big.Int)
		power3.Exp(bigTwo, big.NewInt(int64(i)), nil)
		power4.Exp(bigTwo, power3, phiN)
		y.Mul(y, power4)
		y.Add(y, alpha)
		y.Mod(y, phiN)

		inverseC := new(big.Int).Set(c)
		inverseC.Sub(big.NewInt(0), inverseC)

		result1 := new(big.Int).Set(g)
		result1.Exp(result1, y, groupP.N)
		result2 := new(big.Int).Set(mArray[i])
		result2.Exp(result2, inverseC, groupP.N)
		result1.Mul(result1, result2)
		result1.Mod(result1, groupP.N)

		result3 := new(big.Int).Set(mArray[i])
		result3.Exp(result3, y, groupP.N)
		result4 := new(big.Int).Set(mArray[i+1])
		result4.Exp(result4, inverseC, groupP.N)
		result3.Mul(result3, result4)
		result3.Mod(result3, groupP.N)

		if result1.Cmp(z) != 0 {
			panic(fmt.Errorf("===>[ERROR from GeneratePublicParameter]test1 error"))
		}
		if result3.Cmp(w) != 0 {
			panic(fmt.Errorf("===>[ERROR from GeneratePublicParameter]test2 error"))
		}

		localProofSet := [3]string{z.String(), w.String(), y.String()}
		proofSet = append(proofSet, localProofSet)
	}

	fmt.Println("===>[Setup]pass all Proof tests!")

	return g, mArray, proofSet
}

// get timeParameter T from config file
func GetT() int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetG]Read config file failed:%s", err))
	}

	t := configViper.GetInt("timeParameter")

	return t
}

// get groupLength L from config file
func GetL() int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetG]Read config file failed:%s", err))
	}

	l := configViper.GetInt("groupLength")

	return l
}

// get g from config file
func GetG() *big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetG]Read config file failed:%s", err))
	}

	g := new(big.Int)
	g.SetString(configViper.GetString("g"), 10)

	return g
}

// get N from config file
func GetN() *big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetN]Read config file failed:%s", err))
	}

	N := new(big.Int)
	N.SetString(configViper.GetString("N"), 10)

	return N
}

// get N from config file
func GetMArray() []*big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetMArray]Read config file failed:%s", err))
	}

	var mArray []*big.Int
	mArrayString := configViper.GetStringSlice("mArray")
	for _, m := range mArrayString {
		mBigint, _ := new(big.Int).SetString(m, 10)
		mArray = append(mArray, mBigint)
	}

	return mArray
}
