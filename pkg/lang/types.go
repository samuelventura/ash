package ash

import (
	"math/big"
)

type edtQuantity struct {
	number *big.Rat
	unit   string
}

func newEdtNumber(value string) *big.Rat {
	r := new(big.Rat)
	r.SetString(value)
	return r
}

func newEdtQuantity(value string, unit string) *edtQuantity {
	edt := new(edtQuantity)
	edt.number = newEdtNumber(value)
	edt.unit = unit
	return edt
}
