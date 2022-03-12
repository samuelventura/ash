package ash

import "fmt"

const (
	errorIndent = iota
	errorParse
)

type errorDo struct {
	tid  int
	row  int
	col  int
	desc string
}

func (do *errorDo) Error() string {
	return fmt.Sprintf("%d:%d:%s", do.row, do.col, do.desc)
}
