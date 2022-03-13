package ash

func newEngine() *engineDo {
	do := new(engineDo)
	do.lexer = defaultLexer()
	return do
}

type engineDo struct {
	lexer  func(string) ([]*tokenDo, int)
	parser func([]*tokenDo) *clauseDo
}

func (do *engineDo) executeString(code string) interface{} {
	return do.executeCode(newFile("", code, true))
}

func (do *engineDo) executeCode(cdo *fileDo) interface{} {

	return nil
}
