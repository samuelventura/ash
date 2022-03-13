package ash

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
