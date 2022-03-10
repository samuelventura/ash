package ash

type engine struct {
}

func (e *engine) assert_equal(code string, value interface{}) {

}

func (e *engine) execute(code string) interface{} {
	return nil
}

func Compile(script string) *engine {
	return &engine{}
}
