package tc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
)

var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

type GroupParameter struct {
	N      *big.Int
	Primes []*big.Int
}

type Configurations struct {
	G      string
	N      string
	MArray []string
}

func GenerateGroupParameter(random io.Reader, bits int) (*GroupParameter, error) {
	MaybeReadByte(random)
	groupP := new(GroupParameter)

	if bits < 64 {
		primeLimit := float64(uint64(1) << uint(bits/2))
		// pi approximates the number of primes less than primeLimit
		pi := primeLimit / (math.Log(primeLimit) - 1)
		// Generated primes start wth 11 (in binary) so we can only
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
		// wll give good results.
		for i := 0; i < 2; i++ {
			var err error
			primes[i], err = rand.Prime(random, todo/(2-i))
			if err != nil {
				return nil, err
			}
			todo -= primes[i].BitLen()
		}

		// Make sure that primes is pairwse unequal.
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

func GeneratePublicParameter(groupP *GroupParameter, bits int, k int) (*big.Int, []*big.Int, [][3]string) {
	var mArray []*big.Int
	primesUnder128 := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127}

	h, err := rand.Int(rand.Reader, groupP.N)
	if err != nil {
		fmt.Println("[Setup]generate h wrong", err)
	}
	// fmt.Println("[Setup]H is", h)

	phiN := new(big.Int).Set(bigOne)
	for _, prime := range groupP.Primes {
		primeSubOne := new(big.Int)
		primeSubOne.Sub(prime, bigOne)
		// fmt.Println(primeSubOne)
		phiN.Mul(phiN, primeSubOne)
	}
	// fmt.Println("[Setup]Phi N is", phiN)

	power := new(big.Int).Set(bigOne)
	for _, p := range primesUnder128 {
		q := new(big.Int)
		q.Exp(big.NewInt(p), big.NewInt(int64(bits)), phiN)
		power.Mul(power, q)
		power.Mod(power, phiN)
	}
	g := new(big.Int)
	g.Exp(h, power, groupP.N)
	// fmt.Println("[Setup]G is", g)

	for i := 0; i <= k; i++ {
		// fmt.Println("Number", i)
		power1 := new(big.Int)
		power2 := new(big.Int)
		m := new(big.Int)
		power1.Exp(bigTwo, big.NewInt(int64(i)), nil)
		power2.Exp(bigTwo, power1, phiN)
		m.Exp(g, power2, groupP.N)
		// fmt.Println("Power is", power2)
		// fmt.Println("M is", m)
		mArray = append(mArray, m)

		if i == k {
			mFirst := new(big.Int)
			power1.Sub(power1, bigOne)
			power2.Exp(bigTwo, power1, phiN)
			mFirst.Exp(g, power2, groupP.N)
			mArray = append(mArray, mFirst)
		}
	}
	// fmt.Println("[Setup]Length of m array is", len(mArray))

	var proofSet [][3]string
	for i := 0; i < k; i++ {
		alpha, _ := rand.Int(rand.Reader, phiN)
		z := new(big.Int)
		z.Exp(g, alpha, groupP.N)
		w := new(big.Int)
		w.Exp(mArray[i], alpha, groupP.N)

		gHash := new(big.Int).SetBytes(Digest((g)))
		nHash := new(big.Int).SetBytes(Digest(groupP.N))
		biHash := new(big.Int).SetBytes(Digest(mArray[i]))
		biPlusHash := new(big.Int).SetBytes(Digest(mArray[i+1]))
		zHash := new(big.Int).SetBytes(Digest(z))
		wHash := new(big.Int).SetBytes(Digest(w))

		c := big.NewInt(0)
		c.Xor(c, gHash)
		c.Xor(c, nHash)
		c.Xor(c, biHash)
		c.Xor(c, biPlusHash)
		c.Xor(c, zHash)
		c.Xor(c, wHash)
		// fmt.Println("random number is", c)

		y := new(big.Int).Set(c)
		power3 := new(big.Int)
		power4 := new(big.Int)
		power3.Exp(bigTwo, big.NewInt(int64(i)), nil)
		power4.Exp(bigTwo, power3, phiN)
		y.Mul(y, power4)
		y.Add(y, alpha)
		y.Mod(y, phiN)
		// fmt.Println("commiter proof", y)

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

		// fmt.Println(result1)
		// fmt.Println(result1.Cmp(z))
		// fmt.Println(result3)
		// fmt.Println(result3.Cmp(w))
		if result1.Cmp(z) != 0 {
			fmt.Println("test1 error")
		}
		if result3.Cmp(w) != 0 {
			fmt.Println("test2 error")
		}

		localProofSet := [3]string{z.String(), w.String(), y.String()}
		proofSet = append(proofSet, localProofSet)
	}

	fmt.Println("[Setup]pass all tests!")
	return g, mArray, proofSet
}
