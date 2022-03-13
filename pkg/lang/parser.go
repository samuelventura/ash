package ash

type contextDo struct {
	parent *contextDo
	named  map[string]interface{}
	last   interface{}
}

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
		tokens, pos := lexer(line)
		//incomplete lexing
		if pos != len(line) {
			edo := new(errorDo)
			edo.tid = errorParse
			edo.row = ldo.number
			edo.col = ldo.indent + pos
			edo.desc = "invalid token"
			return nil, edo
		}
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

func scanOne(tid int) func([]*tokenDo) int {
	return func(tokens []*tokenDo) int {
		if hasPrefix(tokens, []int{tid}) {
			return 1
		}
		return 0
	}
}

// func scanNumber(tokens []*tokenDo) int {
// 	return scanOne(tokenNumber)(tokens)
// }

// func scanQuantity(tokens []*tokenDo) int {
// 	return scanOne(tokenQuantity)(tokens)
// }

// func parseLiteral(tokens []*tokenDo) func(ctx *contextDo) {

// }

// func parseReference(tokens []*tokenDo) func(ctx *contextDo) {

// }

func parseJoin(parsers ...func(tokens []*tokenDo) *clauseDo) *clauseDo {

	return nil
}

func parseExpression(tokens []*tokenDo) *clauseDo {

	return nil
}
