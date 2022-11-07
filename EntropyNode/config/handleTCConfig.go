package config

import (
	"fmt"
	"math/big"

	"github.com/spf13/viper"
)

// get g from config
func GetG() *big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	g := new(big.Int)
	g.SetString(configViper.GetString("g"), 10)

	return g
}

// get N from config
func GetN() *big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	N := new(big.Int)
	N.SetString(configViper.GetString("N"), 10)

	return N
}

// get N from config
func GetMArray() []*big.Int {
	// set config file
	configViper := viper.New()
	configViper.SetConfigFile("../TC.yml")

	if err := configViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var mArray []*big.Int
	mArrayString := configViper.GetStringSlice("mArray")
	for _, m := range mArrayString {
		mBigint, _ := new(big.Int).SetString(m, 10)
		mArray = append(mArray, mBigint)
		// fmt.Println(mBigint)
	}
	fmt.Println("Length of m array is", len(mArray))

	return mArray
}
