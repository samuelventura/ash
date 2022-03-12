package ash

type tokenDo struct {
	tid  int
	text string
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
	tokenComment
)

func defaultLexer() func(string) ([]*tokenDo, int) {
	return lexJoin(
		lexComment,
		lexSpace,
		lexName,
		lexQuantity,
		lexNumber,
		lexPrefix("++", tokenPlusPlus),
		lexPrefix("--", tokenMinusMinus),
		lexPrefix("+=", tokenPlusEqual),
		lexPrefix("-=", tokenMinusEqual),
		lexPrefix(":", tokenColon),
		lexPrefix(".", tokenDot),
		lexPrefix("+", tokenPlus),
		lexPrefix("-", tokenMinus),
		lexPrefix("=", tokenEqual),
	)
}

func lexJoin(tokenizers ...func(string) *tokenDo) func(string) ([]*tokenDo, int) {
	return func(line string) ([]*tokenDo, int) {
		pos := 0
		list := new(listDo)
		for pos < len(line) {
			size := 0
			for _, tokenizer := range tokenizers {
				token := tokenizer(line[pos:])
				if token != nil {
					list.append(token)
					size = len(token.text)
					break
				}
			}
			if size > 0 {
				pos += size
				continue
			}
			break
		}
		tokens := make([]*tokenDo, 0, list.length)
		list.each(func(value interface{}) {
			token := value.(*tokenDo)
			tokens = append(tokens, token)
		})
		return tokens, pos
	}
}

func lexPrefix(prefix string, id int) func(line string) *tokenDo {
	return func(line string) *tokenDo {
		size := scanPrefix(prefix)(line)
		if size > 0 {
			return &tokenDo{id, line[:size]}
		}
		return nil
	}
}

func lexScanner(line string, id int, scanner func(line string) int) *tokenDo {
	size := scanner(line)
	if size > 0 {
		return &tokenDo{id, line[:size]}
	}
	return nil
}

func lexComment(line string) *tokenDo {
	return lexScanner(line, tokenComment, scanComment)
}

func lexSpace(line string) *tokenDo {
	return lexScanner(line, tokenSpace, scanSpaces)
}

func lexName(line string) *tokenDo {
	return lexScanner(line, tokenName, scanSome(
		scanAlphaUnders, scanAlphaUnderDigits))
}

func lexQuantity(line string) *tokenDo {
	return lexScanner(line, tokenQuantity, scanAll(
		scanSome(scanDigits, scanAll(scanPrefix("."), scanDigits)),
		scanAlphas,
	))
}

func lexNumber(line string) *tokenDo {
	return lexScanner(line, tokenNumber, scanSome(
		scanDigits,
		scanAll(scanPrefix("."), scanDigits),
	))
}
