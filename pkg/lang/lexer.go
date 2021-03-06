package ash

type tokenDo struct {
	tid  int
	text string
}

const (
	tokenSpace = iota
	tokenComment
	tokenName
	tokenNumber
	tokenPlusPlus
	tokenMinusMinus
	tokenPlusEqual
	tokenMinusEqual
	tokenOpen
	tokenClose
	tokenColon
	tokenDot
	tokenPlus
	tokenMinus
	tokenEqual
)

func buildLexer() func(string) ([]*tokenDo, int) {
	return lexJoin(
		lexComment,
		lexSpace,
		lexName,
		lexNumber,
		lexPrefix("++", tokenPlusPlus),
		lexPrefix("--", tokenMinusMinus),
		lexPrefix("+=", tokenPlusEqual),
		lexPrefix("-=", tokenMinusEqual),
		lexPrefix("(", tokenOpen),
		lexPrefix(")", tokenClose),
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

func lexNumber(line string) *tokenDo {
	return lexScanner(line, tokenNumber, scanOne(
		scanAll(scanPrefix("0b"), scanDigits),
		scanAll(scanPrefix("0o"), scanDigits),
		scanAll(scanPrefix("0x"), scanDigits),
		scanSome(scanDigits, scanAll(scanPrefix("."), scanDigits))))
}
