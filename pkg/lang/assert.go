package ash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newAssert(t *testing.T) *assertDo {
	do := new(assertDo)
	do.t = t
	return do
}

type assertDo struct {
	t *testing.T
}

func (do *assertDo) equal(value interface{}, expected interface{}) {
	assert.Equal(do.t, expected, value)
}

func (do *assertDo) executeEqual(engine *engineDo, code string, expected interface{}) {
	value := engine.executeString(code)
	assert.Equal(do.t, expected, value)
}
