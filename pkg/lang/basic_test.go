package ash

import (
	"testing"
)

func TestNonCodeLinesAreIgnored(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "", nil)
	assert.executeEqual(engine, "#", nil)
	assert.executeEqual(engine, "\n#", nil)
	assert.executeEqual(engine, "12\t#\t\n\t#\t", newDtNumber("12"))
}

func TestLastValueIsReset(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "12", newDtNumber("12"))
	assert.executeEqual(engine, "", nil)
}

func TestNumberLiterals(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "0x12", newDtNumber("18"))
	assert.executeEqual(engine, "0b101", newDtNumber("5"))
	assert.executeEqual(engine, "0o12", newDtNumber("10"))
	assert.executeEqual(engine, "12", newDtNumber("12"))
	assert.executeEqual(engine, "1.2", newDtNumber("1.2"))
}

func TestBasicExpressions(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "12", newDtNumber("12"))
	assert.executeEqual(engine, "12ms", newDtQuantity("12", "ms"))
	// assert.executeEqual(engine, "a=12", newDtNumber("12"))
	// assert.executeEqual(engine, "a", newDtNumber("12"))
}

// func TestBasicAssigments(t *testing.T) {
// 	assert := newAssert(t)
// 	engine := newEngine()
// 	assert.executeEqual(engine, "a=1", newEdtNumber("1"))
// 	assert.executeEqual(engine, "b=1ms", newEdtQuantity("1", "ms"))
// 	assert.executeEqual(engine, "a", newEdtNumber("1"))
// 	assert.executeEqual(engine, "b", newEdtQuantity("1", "ms"))
// }

// func TestBasicTree(t *testing.T) {
// 	assert := newAssert(t)
// 	code := newCode(`
// 	tree global:
// 		count = 0
// 		increment:
// 			count++
// 		add amount:
// 			count += amount
// 	`)
// 	fmt.Println(code.toString())
// 	assert.equal(0, len(code.errors))
// 	engine := newEngine()
// 	engine.executeCode(code)
// 	assert.executeEqual(engine, "global.count", 0)
// 	engine.executeString("global.increment")
// 	assert.executeEqual(engine, "global.count", 1)
// 	engine.executeString("global.add 2")
// 	assert.executeEqual(engine, "global.count", 3)
// }
