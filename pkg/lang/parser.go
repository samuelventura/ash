package ash

import (
	"strings"
)

type clauseDo struct {
	oper   string
	args   []string
	body   []*clauseDo
	bodied bool
}

func parse(indent int, lines []*lineDo, lexer func(string) ([]*tokenDo, int), parser func([]*tokenDo) *clauseDo) ([]*clauseDo, error) {
	//check identation
	current := indent
	for _, ldo := range lines {
		step := ldo.indent - current
		//more than one step at a time
		if step > 1 {
			edo := new(errorDo)
			edo.tid = errorIndent
			edo.row = ldo.number
			edo.col = ldo.indent
			edo.desc = "invalid indent"
			return nil, edo
		}
		current = ldo.indent
	}
	list := new(listDo)
	iter := lineIterator(lines)
	for !iter.done() {
		ldo := iter.pop()
		//must be at current indent
		if ldo.indent != indent {
			edo := new(errorDo)
			edo.tid = errorParse
			edo.row = ldo.number
			edo.col = ldo.indent
			edo.desc = "invalid indent"
			return nil, edo
		}
		line := ldo.text[ldo.indent:]
		//trim trailing spaces
		line = strings.TrimSpace(line)
		tokens, pos := lexer(line)
		//incomplete lexing
		invalid := pos != len(line)
		invalid = invalid || len(tokens) == 0
		if invalid {
			edo := new(errorDo)
			edo.tid = errorParse
			edo.row = ldo.number
			edo.col = ldo.indent + pos
			edo.desc = "invalid token"
			return nil, edo
		}
		//trim trailing tokens
		pos = len(tokens)
		for pos > 0 {
			token := tokens[pos-1]
			valid := token.tid != tokenSpace
			valid = valid && token.tid != tokenComment
			if valid {
				break
			}
			pos--
		}
		tokens = tokens[:pos]
		clause := parser(tokens)
		//no clause found
		if clause == nil {
			edo := new(errorDo)
			edo.tid = errorParse
			edo.row = ldo.number
			edo.col = ldo.indent
			edo.desc = "invalid clause"
			return nil, edo
		}
		if clause.bodied {
			nlist := new(listDo)
			for !iter.done() {
				ndo := iter.peek()
				if ndo.indent > indent {
					nlist.append(ndo)
					iter.pop()
					continue
				}
				break
			}
			if nlist.length == 0 {
				edo := new(errorDo)
				edo.tid = errorParse
				edo.row = ldo.number
				edo.col = ldo.indent
				edo.desc = "empty body"
				return nil, edo
			}
			nlines := make([]*lineDo, 0, nlist.length)
			nlist.each(func(value interface{}) {
				line := value.(*lineDo)
				nlines = append(nlines, line)
			})
			body, err := parse(indent+1, nlines, lexer, parser)
			if err != nil {
				return nil, err
			}
			clause.body = body
		}
		list.append(clause)
	}
	clauses := make([]*clauseDo, 0, list.length)
	list.each(func(value interface{}) {
		clause := value.(*clauseDo)
		clauses = append(clauses, clause)
	})
	return clauses, nil
}

func hasPrefix(tokens []*tokenDo, tids []int) bool {
	tokenc := len(tokens)
	for i, tid := range tids {
		if i >= tokenc {
			return false
		}
		if tid != tokens[i].tid {
			return false
		}
	}
	return true
}

func peekOr(peekers ...func([]*tokenDo) int) func([]*tokenDo) int {
	return func(tokens []*tokenDo) int {
		for _, peeker := range peekers {
			size := peeker(tokens)
			if size > 0 {
				return size
			}
		}
		return 0
	}
}

func peekAnd(peekers ...func([]*tokenDo) int) func([]*tokenDo) int {
	return func(tokens []*tokenDo) int {
		pos := 0
		for _, peeker := range peekers {
			size := peeker(tokens[pos:])
			if size > 0 {
				pos += size
				continue
			}
			return 0
		}
		return pos
	}
}

func peekTids(tids ...int) func([]*tokenDo) int {
	return func(tokens []*tokenDo) int {
		if hasPrefix(tokens, tids) {
			return len(tids)
		}
		return 0
	}
}

func parseJoin(parsers ...func([]*tokenDo) *clauseDo) func([]*tokenDo) *clauseDo {
	return func(tokens []*tokenDo) *clauseDo {
		for _, parser := range parsers {
			clause := parser(tokens)
			if clause != nil {
				return clause
			}
		}
		return nil
	}
}

func peek(tokens []*tokenDo, peeker func([]*tokenDo) int) bool {
	return peeker(tokens) == len(tokens)
}

func parserQuantity(tokens []*tokenDo) *clauseDo {
	if peek(tokens, peekTids(tokenNumber, tokenName)) {
		cdo := new(clauseDo)
		cdo.oper = "lq"
		cdo.args = []string{tokens[0].text, tokens[1].text}
		return cdo
	}
	return nil
}

func parserNumber(tokens []*tokenDo) *clauseDo {
	if peek(tokens, peekTids(tokenNumber)) {
		cdo := new(clauseDo)
		cdo.oper = "ln"
		cdo.args = []string{tokens[0].text}
		return cdo
	}
	return nil
}

func buildParser() func(tokens []*tokenDo) *clauseDo {
	return parseJoin(
		parserQuantity,
		parserNumber,
	)
}
