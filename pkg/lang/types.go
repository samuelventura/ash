package ash

import (
	"fmt"
	"math/big"
)

type dtNumber struct {
	value *big.Rat
}

func (edt *dtNumber) String() string {
	return edt.value.FloatString(8)
}

type dtQuantity struct {
	number *dtNumber
	unit   string
}

func (edt *dtQuantity) String() string {
	return fmt.Sprintf("%v%s", edt.number, edt.unit)
}

func newDtNumber(args ...string) *dtNumber {
	edt := new(dtNumber)
	edt.value = new(big.Rat)
	edt.value.SetString(args[0])
	return edt
}

func newDtQuantity(args ...string) *dtQuantity {
	edt := new(dtQuantity)
	edt.number = newDtNumber(args[0])
	edt.unit = args[1]
	return edt
}
