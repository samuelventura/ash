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
	do.context.parent = nil
	do.context.named = make(map[string]interface{})
	do.parser = buildParser()
	do.lexer = buildLexer()
	return do
}

func (do *engineDo) get(name string) interface{} {
	return do.context.named[name]
}

func (do *engineDo) set(name string, value interface{}) {
	do.context.named[name] = value
}

func (do *engineDo) executeString(code string) interface{} {
	return do.executeCode(newFile("", code, true))
}

func (do *engineDo) executeCode(cdo *fileDo) interface{} {
	clauses, err := parse(0, cdo.lines, do.lexer, do.parser)
	if err != nil {
		panic(err)
	}
	do.context.last = nil
	for _, clause := range clauses {
		do.context.last = clause.exec(do.context)
	}
	return do.context.last
}
