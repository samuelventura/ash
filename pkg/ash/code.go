package ash

import (
	"strings"
	"unicode"
)

const (
	lineEof = iota
	lineEmpty
	lineComment
	lineCode
)

type lineDo struct {
	tid    int
	number int
	indent int
	text   string
	tokens []*tokenDo
	apply  func(ctx *contextDo)
}

const (
	errorIndent = iota
	errorParse
)

type errorDo struct {
	tid    int
	line   *lineDo
	column int
	desc   string
}

type codeDo struct {
	code   string
	indent int
	lines  []*lineDo
}

func (do *codeDo) toString() string {
	b := new(strings.Builder)
	for _, ldo := range do.lines {
		b.WriteString(ldo.text)
		b.WriteRune('\n')
	}
	return b.String()
}

//code indent analyzer removes common prefix
//for testing purposes, strict otherwise
//comments must be equally indented as well
//spaces and tabs count equally as one
//leave format fixing to the code editor
func newCode(code string) *codeDo {
	lines := strings.Split(code, "\n")
	cdo := new(codeDo)
	cdo.code = code
	cdo.lines = make([]*lineDo, 0, len(lines))
	cdo.indent = -1
	for i, line := range lines {
		first_j := -1
		first_c := '\n'
		for j, c := range line {
			if !unicode.IsSpace(c) {
				first_j = j
				first_c = c
				break
			}
		}
		if first_j >= 0 {
			if cdo.indent < 0 || first_j < cdo.indent {
				cdo.indent = first_j
			}
		}
		ldo := new(lineDo)
		ldo.number = i
		ldo.indent = first_j
		ldo.text = line
		switch first_c {
		case '#':
			ldo.tid = lineComment
		case '\n':
			ldo.tid = lineEmpty
		default:
			ldo.tid = lineCode
		}
		cdo.lines = append(cdo.lines, ldo)
	}
	//shift identation left
	if cdo.indent > 0 {
		for _, ldo := range cdo.lines {
			ldo.indent -= cdo.indent
			ldo.text = ldo.text[cdo.indent:]
		}
	}
	return cdo
}
