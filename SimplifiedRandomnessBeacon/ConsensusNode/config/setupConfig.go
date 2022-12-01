package config

import (
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

// config file Setup
// assign curve,previous output,VRF type,Time Commitment Type,F Function Type
func SetupConfig() {
	// lock file
	f, err := os.Open("../lock")
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Open lock failed:%s", err))
	}
	// share lock, concurrently read
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Add share lock failed:%s", err))
	}

	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/config.yml")
	outputViper := viper.New()
	outputViper.SetConfigFile("../Configuration/output.yml")
	tcViper := viper.New()
	tcViper.SetConfigFile("../Configuration/TC.yml")

	// read config and keep origin settings
	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Read config file failed:%s", err))
	}
	if err := outputViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Read config file failed:%s", err))
	}
	if err := tcViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Read config file failed:%s", err))
	}

	// first time setup
	if !configViper.GetBool("Running") {
		fmt.Println("===>[Setup]First time Setup")

		// generate elliptic curve
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]setup conf curve(private Key) failed, err:%s", err))
		}
		marshalledKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]setup conf curve(marshalled Key) failed, err:%s", err))
		}
		configViper.Set("EllipticCurve", string(marshalledKey))

		// generate random difficulty
		selectBigInt, _ := rand.Int(rand.Reader, big.NewInt(2))
		selectInt, err := strconv.Atoi(selectBigInt.String())
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]generate random difficulty failed, err:%s", err))
		}
		configViper.Set("Difficulty", selectInt)

		configViper.Set("Running", true)
		// TODO: VRF type,Time Commitment Type,F Function Type

		// generate random init output
		message := []byte("asdkjhdk")
		randomNum := util.Digest(message)
		outputViper.Set("PreviousOutput", string(randomNum))

		// group parameter contains prime numbers p,q and a large number n=p*q of groupLength bits
		groupLength := GetL()
		groupParameter, err := GenerateGroupParameter(rand.Reader, groupLength)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]generate group parameter failed, err:%s", err))
		}
		fmt.Println("===>[Setup]N is", groupParameter.N)
		fmt.Println("===>[Setup]primes are", groupParameter.Primes)
		fmt.Println("===>[Setup]Is prime[0] prime?", groupParameter.Primes[0].ProbablyPrime(20))
		fmt.Println("===>[Setup]Is prime[1] prime?", groupParameter.Primes[1].ProbablyPrime(20))

		// generator g
		timeParameter := GetT()
		g, mArray, proofSet := GeneratePublicParameter(groupParameter, groupLength, timeParameter)
		fmt.Println("===>[Setup]time Parameter T is", timeParameter)
		fmt.Println("===>[Setup]G is", g)
		fmt.Println("===>[Setup]Length of m array is", len(mArray))

		// m array
		var mArrayString []string
		for _, m := range mArray {
			mArrayString = append(mArrayString, m.String())
		}

		tcViper.Set("g", g.String())
		tcViper.Set("n", groupParameter.N.String())
		tcViper.Set("prime0", groupParameter.Primes[0].String())
		tcViper.Set("prime1", groupParameter.Primes[1].String())
		tcViper.Set("mArray", mArrayString)
		tcViper.Set("proofSet", proofSet)

		// write new settings
		if err := configViper.WriteConfig(); err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]write config failed, err:%s", err))
		}
		if err := outputViper.WriteConfig(); err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]write output failed, err:%s", err))
		}
		if err := tcViper.WriteConfig(); err != nil {
			panic(fmt.Errorf("===>[ERROR from SetupConfig]write tc failed, err:%s", err))
		}

		fmt.Println("===>[Setup]Finish first time Setup")
	}

	// unlock file
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]Unlock share lock failed:%s", err))
	}
}
