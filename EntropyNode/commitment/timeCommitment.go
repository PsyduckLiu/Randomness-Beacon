package commitment

import (
	"crypto/rand"
	"math/big"
)

func GenerateTimeCommitment() *big.Int {
	var upper, e = big.NewInt(2), big.NewInt(256)
	upper.Exp(upper, e, nil)

	// generate random number
	tc, err := rand.Int(rand.Reader, upper)
	if err != nil {
		panic(err)
	}

	return tc
}
