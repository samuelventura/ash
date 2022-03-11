package ash

func newEngine() *engineDo {
	do := new(engineDo)
	do.tokenizers = []func(string) *tokenDo{
		tokenizeSpace,
		tokenizeName,
		tokenizeQuantity,
		tokenizeNumber,
		tokenizePrefix("++", tokenPlusPlus),
		tokenizePrefix("--", tokenMinusMinus),
		tokenizePrefix("+=", tokenPlusEqual),
		tokenizePrefix("-=", tokenMinusEqual),
		tokenizePrefix(":", tokenColon),
		tokenizePrefix(".", tokenDot),
		tokenizePrefix("+", tokenPlus),
		tokenizePrefix("-", tokenMinus),
		tokenizePrefix("=", tokenEqual),
		tokenizeAny,
	}
	return do
}

type engineDo struct {
	tokenizers []func(string) *tokenDo
}

func (do *engineDo) executeString(code string) interface{} {
	return do.executeCode(newCode(code))
}

func (do *engineDo) executeCode(code *codeDo) interface{} {
	return nil
}
