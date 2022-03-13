package ash

import (
	"strings"
	"unicode"
)

const (
	runeEol = iota
	runeSpace
	runeAlpha
	runeDigit
	runeUnder
	runeDot
	runeColon
	runeEqual
	runePlus
	runeMinus
	runeAny
)

func runeFlags(r rune) int {
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
			return runeAlpha
		} else if unicode.IsSpace(r) {
			return runeSpace
		} else {
			return runeAny
		}
	}
}

type runeDo struct {
	r rune //rune
	i int  //index
	f int  //flags
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

func scanComment(line string) int {
	if strings.HasPrefix(line, "#") {
		return len(line)
	} else {
		return 0
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
	pos := 0
	for _, r := range line {
		f := runeFlags(r)
		if !valid(f) {
			break
		}
		pos++
	}
	return pos
}

func scanSpaces(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeSpace
	})
}

func scanAlphas(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha
	})
}

func scanAlphaUnders(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha || id == runeUnder
	})
}

func scanAlphaUnderDigits(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeAlpha || id == runeUnder || id == runeDigit
	})
}

func scanDigits(line string) int {
	return scanValid(line, func(id int) bool {
		return id == runeDigit
	})
}
