package ash

type contextDo struct {
	parent *contextDo
	named  map[string]interface{}
	last   interface{}
}

type engineDo struct {
	context *contextDo
	parser  func([]*tokenDo) *clauseDo
	lexer   func(string) ([]*tokenDo, int)
}

func newEngine() *engineDo {
	do := new(engineDo)
	do.context = new(contextDo)
	do.parser = buildParser()
	do.lexer = buildLexer()
	return do
}

func (do *engineDo) executeString(code string) interface{} {
	return do.executeCode(newFile("", code, true))
}

func (do *engineDo) executeCode(cdo *fileDo) interface{} {
	clauses, err := parse(0, cdo.lines, do.lexer, do.parser)
	if err != nil {
		panic(err)
	}
	for _, clause := range clauses {
		execute(clause, do.context)
	}
	return do.context.last
}

func execute(cdo *clauseDo, ctx *contextDo) {
	switch cdo.oper {
	case "ln":
		ctx.last = newDtNumber(cdo.args...)
	case "lq":
		ctx.last = newDtQuantity(cdo.args...)
	}
}
