package ash

type tokenDo struct {
	tid   int
	value string
}

const (
	tokenEol = iota
	tokenSpace
	tokenName
	tokenQuantity
	tokenNumber
	tokenPlusPlus
	tokenMinusMinus
	tokenPlusEqual
	tokenMinusEqual
	tokenColon
	tokenDot
	tokenPlus
	tokenMinus
	tokenEqual
	tokenError
)

func defaultTokenizer() func(string) []*tokenDo {
	return tokenizeJoin(
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
	)
}

func tokenizeJoin(tokenizers ...func(string) *tokenDo) func(string) []*tokenDo {
	return func(line string) []*tokenDo {
		pos := 0
		tokens := make([]*tokenDo, 0, 4)
		for pos < len(line) {
			curr := pos
			for _, tokenizer := range tokenizers {
				token := tokenizer(line[pos:])
				if token != nil {
					tokens = append(tokens, token)
					pos += len(token.value)
					break
				}
			}
			if curr == pos {
				token := &tokenDo{tokenError, line[pos:]}
				tokens = append(tokens, token)
				pos = len(line)
			}
		}
		return tokens
	}
}

func tokenizePrefix(prefix string, id int) func(line string) *tokenDo {
	return func(line string) *tokenDo {
		size := scanPrefix(prefix)(line)
		if size > 0 {
			return &tokenDo{id, line[:size]}
		}
		return nil
	}
}

func tokenizeScanner(line string, id int, scanner func(line string) int) *tokenDo {
	size := scanner(line)
	if size > 0 {
		return &tokenDo{id, line[:size]}
	}
	return nil
}

func tokenizeSpace(line string) *tokenDo {
	return tokenizeScanner(line, tokenSpace, scanSpace)
}

func tokenizeName(line string) *tokenDo {
	return tokenizeScanner(line, tokenName, scanSome(
		scanAlphaUnder, scanAlphaUnderDigit))
}

func tokenizeQuantity(line string) *tokenDo {
	return tokenizeScanner(line, tokenQuantity, scanAll(
		scanSome(scanDigits, scanAll(scanPrefix("."), scanDigits)),
		scanAlpha,
	))
}

func tokenizeNumber(line string) *tokenDo {
	return tokenizeScanner(line, tokenNumber, scanSome(
		scanDigits,
		scanAll(scanPrefix("."), scanDigits),
	))
}
