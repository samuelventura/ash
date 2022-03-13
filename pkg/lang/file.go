package ash

import (
	"strings"
	"unicode"
)

type lineDo struct {
	number int
	indent int
	text   string
	file   *fileDo
}

type fileDo struct {
	name  string    //filename
	text  string    //original text
	lines []*lineDo //parsed lines
}

func (do *fileDo) toString() string {
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
func newFile(name string, text string, fixIndent bool) *fileDo {
	lines := strings.Split(text, "\n")
	fdo := new(fileDo)
	fdo.text = text
	fdo.name = name
	indent := -1
	list := new(listDo)
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
			if indent < 0 || first_j < indent {
				indent = first_j
			}
		}
		ldo := new(lineDo)
		ldo.file = fdo
		ldo.number = i
		ldo.indent = first_j
		ldo.text = line
		switch first_c {
		case '#':
			//ignore
		case '\n':
			//ignore
		default:
			list.append(ldo)
		}
	}
	fdo.lines = make([]*lineDo, 0, list.length)
	list.each(func(value interface{}) {
		line := value.(*lineDo)
		fdo.lines = append(fdo.lines, line)
	})
	if fixIndent && indent > 0 {
		//shift identation left
		for _, ldo := range fdo.lines {
			ldo.indent -= indent
			ldo.text = ldo.text[indent:]
		}
	}
	return fdo
}
