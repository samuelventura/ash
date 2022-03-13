package ash

import (
	"strings"
)

type clauseDo struct {
	bodied bool
	line   *lineDo
	body   []*clauseDo
	tag    interface{}
	exec   func(ctx *contextDo) interface{}
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
		clause.line = ldo
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

func peekPrefix(tids ...int) func(tokens []*tokenDo) int {
	return func(tokens []*tokenDo) int {
		tokenc := len(tokens)
		for i, tid := range tids {
			if i >= tokenc {
				return 0
			}
			if tid != tokens[i].tid {
				return 0
			}
		}
		return len(tids)
	}
}

func peekPrefixEx(tids []int, opts []bool) func(tokens []*tokenDo) int {
	return func(tokens []*tokenDo) int {
		off := 0
		tokenc := len(tokens)
		for i, tid := range tids {
			if i-off >= tokenc {
				return 0
			}
			opt := opts[i]
			if tid != tokens[i-off].tid {
				if opt {
					off++
				} else {
					return 0
				}
			}
		}
		return len(tids) - off
	}
}

func peekOneMany(one func([]*tokenDo) int, many func([]*tokenDo) int) func(tokens []*tokenDo) int {
	return func(tokens []*tokenDo) int {
		pos := one(tokens)
		if pos > 0 {
			for {
				size := many(tokens[pos:])
				if size > 0 {
					pos += size
					continue
				}
				break
			}
		}
		return pos
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

func peekReference(tokens []*tokenDo) int {
	one := peekPrefix(tokenName)
	many := peekPrefix(tokenDot, tokenName)
	return peekOneMany(one, many)(tokens)
}

func parserQuantity(tokens []*tokenDo) *clauseDo {
	if peek(tokens, peekPrefix(tokenNumber, tokenName)) {
		cdo := new(clauseDo)
		cdo.tag = newDtQuantity(tokens[0].text, tokens[1].text)
		cdo.exec = func(ctx *contextDo) interface{} {
			return cdo.tag.(*dtQuantity)
		}
		return cdo
	}
	return nil
}

func parserNumber(tokens []*tokenDo) *clauseDo {
	if peek(tokens, peekPrefix(tokenNumber)) {
		cdo := new(clauseDo)
		cdo.tag = newDtNumber(tokens[0].text)
		cdo.exec = func(ctx *contextDo) interface{} {
			return cdo.tag.(*dtNumber)
		}
		return cdo
	}
	return nil
}

func parserReference(tokens []*tokenDo) *clauseDo {
	pos := peekReference(tokens)
	if pos > 0 {
		refc := 1 + (pos-1)/2
		refn := make([]string, refc)
		for i := range refn {
			refn[i] = tokens[2*i].text
		}
		cdo := new(clauseDo)
		cdo.tag = refn
		cdo.exec = func(ctx *contextDo) interface{} {
			refn := cdo.tag.([]string)
			return ctx.named[refn[0]]
		}
		return cdo
	}
	return nil
}

func parserAssigment(tokens []*tokenDo) *clauseDo {
	pos := peekReference(tokens)
	if pos > 0 {
		refc := 1 + (pos-1)/2
		tids := []int{tokenSpace, tokenEqual, tokenSpace}
		opts := []bool{true, false, true}
		size := peekPrefixEx(tids, opts)(tokens[pos:])
		if size > 0 {
			pos += size
			exp := parserExpression(tokens[pos:])
			if exp != nil {
				cl := new(clAssigment)
				cl.expression = exp
				cl.names = make([]string, refc)
				for i := range cl.names {
					cl.names[i] = tokens[2*i].text
				}
				cdo := new(clauseDo)
				cdo.tag = cl
				cdo.exec = func(ctx *contextDo) interface{} {
					cl := cdo.tag.(*clAssigment)
					val := cl.expression.exec(ctx)
					ctx.named[cl.names[0]] = val
					return val
				}
				return cdo
			}
		}
	}
	return nil
}

func parserExpression(tokens []*tokenDo) *clauseDo {
	return parseJoin(
		parserReference,
		parserQuantity,
		parserNumber,
	)(tokens)
}

func buildParser() func(tokens []*tokenDo) *clauseDo {
	return parseJoin(
		parserAssigment,
		parserExpression,
	)
}
