package config

import (
	"fmt"
	"math/big"

	"github.com/spf13/viper"
)

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

// get primes from config file
func GetPrimes() []*big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../Configuration/TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("===>[ERROR from GetPrimes]Read config file failed:%s", err))
	}

	var primes []*big.Int
	prime0 := new(big.Int)
	prime0.SetString(configViper.GetString("prime0"), 10)
	primes = append(primes, prime0)
	prime1 := new(big.Int)
	prime1.SetString(configViper.GetString("prime1"), 10)
	primes = append(primes, prime1)

	return primes
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
