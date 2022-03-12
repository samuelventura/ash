package ash

type contextDo struct {
	parent *contextDo
	named  map[string]interface{}
	last   interface{}
}

func parse(cdo *codeDo, tokenizer func(string) []*tokenDo) func(ctx *contextDo) {
	//check identation
	errors := make([]*errorDo, 0, len(cdo.lines))
	current := -1
	for _, ldo := range cdo.lines {
		step := ldo.indent - current
		if step > 1 {
			edo := new(errorDo)
			edo.tid = errorIndent
			edo.column = ldo.indent
			edo.line = ldo
			edo.desc = "invalid indent"
			errors = append(errors, edo)
		}
		current = ldo.indent
	}
	if len(errors) > 0 {
		return nil
	}
	for _, ldo := range cdo.lines {
		switch ldo.tid {
		case lineCode:
			ldo.tokens = tokenizer(ldo.text[ldo.indent:])
			if len(ldo.tokens) == 0 {
				edo := new(errorDo)
				edo.tid = errorParse
				edo.column = 0
				edo.line = ldo
				edo.desc = "invalid token"
				errors = append(errors, edo)
			} else {
				pos := 0
				for _, token := range ldo.tokens {
					if token.tid == tokenError {
						edo := new(errorDo)
						edo.tid = errorParse
						edo.column = pos
						edo.line = ldo
						edo.desc = "invalid token"
						errors = append(errors, edo)
						break
					}
					pos += len(token.value)
				}
			}
		}
	}
	if len(errors) > 0 {
		return nil
	}
	return nil
}

func parse(tokens []*tokenDo) func(ctx *contextDo) {
	iter := tokenIterator(tokens)
	for !iter.done() {
		ido := iter.next()
		if !valid(ido.tid) {
			return ido.i
		}
	}
	return nil
}

func parseLiteral(tokens []*tokenDo) func(ctx *contextDo) {

}

func parseReference(tokens []*tokenDo) func(ctx *contextDo) {

}

func parseExpression(tokens []*tokenDo) func(ctx *contextDo) {

}
