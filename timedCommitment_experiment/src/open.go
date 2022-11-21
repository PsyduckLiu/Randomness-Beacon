package tc

import (
	"fmt"
	"math/big"
	"tc/config"
)

func ForcedOped(c *big.Int, h *big.Int, rKSubOne *big.Int, rK *big.Int, k int) {
	N := config.GetN()
	// fmt.Println("[Open]Get N is", N)

	var kBig = big.NewInt(int64(k - 1))
	times := new(big.Int)
	times.Exp(bigTwo, kBig, nil)
	times.Sub(times, bigOne)
	fmt.Println("[Open]times is", times)

	var index = big.NewInt(0)
	r := new(big.Int).Set(rKSubOne)
	for index.Cmp(times) == -1 {
		r.Exp(r, bigTwo, N)
		index.Add(index, bigOne)
	}

	// fmt.Println("[Open]After forced open, r is", r)

	msg := new(big.Int)
	msg.Xor(c, r)
	// fmt.Println("[Open]After forced open, msg is", msg)

	rSquare := new(big.Int).Set(r)
	rSquare.Exp(rSquare, bigTwo, N)
	fmt.Println("[Open]Verify", rSquare.Cmp(rK))
}
