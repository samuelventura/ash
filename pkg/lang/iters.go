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
			rdo := new(runeDo)
			rdo.i = pos
			rdo.r = runes[pos]
			rdo.f = runeFlags(rdo.r)
			return rdo
		} else if pos == runec {
			rdo := new(runeDo)
			rdo.i = runec
			rdo.r = '\n'
			rdo.f = runeEol
			return rdo
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
	pos    int
	length int
	lines  []*lineDo
}

func (iter *lineIter) done() bool {
	return iter.pos >= iter.length
}

func (iter *lineIter) peek() *lineDo {
	if iter.pos < iter.length {
		return iter.lines[iter.pos]
	} else {
		return nil
	}
}

func (iter *lineIter) pop() *lineDo {
	defer func() { iter.pos++ }()
	return iter.peek()
}

func lineIterator(lines []*lineDo) *lineIter {
	iter := new(lineIter)
	iter.pos = 0
	iter.lines = lines
	iter.length = len(lines)
	return iter
}
