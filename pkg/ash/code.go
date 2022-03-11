package ash

import (
	"strings"
	"unicode"
)

type lineDo struct {
	number int
	indent int
	text   string
}

type errorDo struct {
	line   *lineDo
	column int
	desc   string
}

type codeDo struct {
	code   string
	count  int
	indent int
	lines  []*lineDo
	errors []*errorDo
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
	cdo.count = 0
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
			if first_c != '#' {
				ldo := new(lineDo)
				ldo.number = i
				ldo.indent = first_j
				ldo.text = line
				cdo.lines = append(cdo.lines, ldo)
				cdo.count++
			}
		}
	}
	//shift identation left
	if cdo.indent > 0 {
		for _, ldo := range cdo.lines {
			ldo.indent -= cdo.indent
			ldo.text = ldo.text[cdo.indent:]
		}
	}
	//check identation
	cdo.errors = make([]*errorDo, 0, len(cdo.lines))
	current := -1
	for _, ldo := range cdo.lines {
		step := ldo.indent - current
		if step > 1 {
			edo := new(errorDo)
			edo.column = ldo.indent
			edo.line = ldo
			edo.desc = "invalid indent"
			cdo.errors = append(cdo.errors, edo)
		}
		current = ldo.indent
	}
	return cdo
}
