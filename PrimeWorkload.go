
package main

import "math/big"

type PrimeWorkload struct {
	Number string;
	IsPrime bool;
}

func (pw *PrimeWorkload) BigInt() *big.Int{
	r := new(big.Int)
	r.SetString(pw.Number,10)
	return r
}
