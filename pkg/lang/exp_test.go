package ash

import (
	"testing"
)

func TestExpressions(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	engine.set("a", newDtNumber("1"))
	engine.set("b", newDtQuantity("2", "ms"))
	assert.executeEqual(engine, "a", newDtNumber("1"))
	assert.executeEqual(engine, "b", newDtQuantity("2", "ms"))
	// assert.executeEqual(engine, "1+2", newDtNumber("3"))
	// assert.executeEqual(engine, "12ms+34ms", newDtQuantity("46", "ms"))
	// assert.executeEqual(engine, "3+a", newDtNumber("4"))
	// assert.executeEqual(engine, "12ms+b", newDtQuantity("14", "ms"))
}
