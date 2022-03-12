package ash

type runeIter struct {
	done func() bool
	next func() *runeDo
}

func runeIterator(line string) *runeIter {
	pos := 0
	runes := []rune(line)
	runec := len(runes)
	iter := new(runeIter)
	iter.done = func() bool { return pos > runec }
	iter.next = func() *runeDo {
		defer func() { pos++ }()
		if pos < runec {
			r := runes[pos]
			return &runeDo{pos, r, runeize(r)}
		} else if pos == runec {
			return &runeDo{pos, '\n', runeEol}
		} else {
			return nil
		}
	}
	return iter
}

type tokenIter struct {
	done func() bool
	next func() *tokenDo
}

func tokenIterator(tokens []*tokenDo) *tokenIter {
	pos := 0
	iter := new(tokenIter)
	tokenc := len(tokens)
	iter.done = func() bool { return pos > tokenc }
	iter.next = func() *tokenDo {
		defer func() { pos++ }()
		if pos < tokenc {
			return tokens[pos]
		} else if pos == tokenc {
			return &tokenDo{tid: tokenEol}
		} else {
			return nil
		}
	}
	return iter
}

type lineIter struct {
	done func() bool
	next func() *lineDo
}

func lineIterator(code *codeDo) *lineIter {
	pos := 0
	iter := new(lineIter)
	linec := len(code.lines)
	iter.done = func() bool { return pos > linec }
	iter.next = func() *lineDo {
		defer func() { pos++ }()
		if pos < linec {
			return code.lines[pos]
		} else if pos == linec {
			return &lineDo{tid: lineEof, number: linec}
		} else {
			return nil
		}
	}
	return iter
}
