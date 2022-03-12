package ash

func newEngine() *engineDo {
	do := new(engineDo)
	do.tokenizer = defaultTokenizer()
	return do
}

type engineDo struct {
	tokenizer func(string) []*tokenDo
}

func (do *engineDo) executeString(code string) interface{} {
	return do.executeCode(newCode(code))
}

func (do *engineDo) executeCode(code *codeDo) interface{} {
	return nil
}
