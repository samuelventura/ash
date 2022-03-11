package ash

import (
	"strings"
	"unicode"
)

const (
	runeAny = iota
	runeSpace
	runeAlpha
	runeDigit
	runeUnder
	runeDot
	runeColon
	runeEqual
	runePlus
	runeMinus
	runeEol
)

func runeize(r rune) int {
	switch r {
	case '.':
		return runeDot
	case ':':
		return runeColon
	case '+':
		return runePlus
	case '-':
		return runeMinus
	case '_':
		return runeUnder
	default:
		if unicode.IsDigit(r) {
			return runeDigit
		} else if unicode.IsLetter(r) {
			return runeDigit
		} else if unicode.IsSpace(r) {
			return runeSpace
		} else {
			return runeAny
		}
	}
}

type runeDo struct {
	i  int
	r  rune
	id int
}

type runeIter struct {
	done func() bool
	next func() *runeDo
}

func runeizer(line string) *runeIter {
	pos := 0
	runes := []rune(line)
	iter := new(runeIter)
	iter.done = func() bool { return pos > len(runes) }
	iter.next = func() *runeDo {
		defer func() { pos++ }()
		if pos < len(runes) {
			r := runes[pos]
			return &runeDo{pos, r, runeize(r)}
		} else if pos == len(runes) {
			return &runeDo{pos, '\n', runeEol}
		} else {
			return nil
		}
	}
	return iter
}

func scanAll(scanners ...func(string) int) func(string) int {
	return func(line string) int {
		pos := 0
		for _, scanner := range scanners {
			size := scanner(line[pos:])
			if size > 0 {
				pos += size
				continue
			}
			return 0
		}
		return pos
	}
}

func scanSome(scanners ...func(string) int) func(string) int {
	return func(line string) int {
		pos := 0
		for _, scanner := range scanners {
			size := scanner(line[pos:])
			if size > 0 {
				pos += size
				continue
			}
			break
		}
		return pos
	}
}

func scanPrefix(prefix string) func(string) int {
	return func(line string) int {
		if strings.HasPrefix(line, prefix) {
			return len(prefix)
		} else {
			return 0
		}
	}
}

func scanValid(line string, valid func(int) bool) int {
	runer := runeizer(line)
	for !runer.done() {
		rdo := runer.next()
		if !valid(rdo.id) {
			return rdo.i
		}
	}
	return 0
}

func scanSpace(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeSpace
	})
}

func scanAlpha(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha
	})
}

func scanAlphaUnder(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha || id == runeUnder
	})
}

func scanAlphaUnderDigit(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha || id == runeUnder || id == runeDigit
	})
}

func scanDigits(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeDigit
	})
}
