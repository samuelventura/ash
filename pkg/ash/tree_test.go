package ash

import (
	"fmt"
	"testing"
)

func TestBasicExpressions(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "1", newEdtNumber("1"))
	assert.executeEqual(engine, "1ms", newEdtQuantity("1", "ms"))
}

func TestBasicTree(t *testing.T) {
	assert := newAssert(t)
	code := newCode(`
	tree global:
		count = 0
		increment:
			count++
		add amount:
			count += amount
	`)
	fmt.Println(code.toString())
	assert.equal(0, len(code.errors))
	engine := newEngine()
	engine.executeCode(code)
	assert.executeEqual(engine, "global.count", 0)
	engine.executeString("global.increment")
	assert.executeEqual(engine, "global.count", 1)
	engine.executeString("global.add 2")
	assert.executeEqual(engine, "global.count", 3)
}
