package ash

import (
	"testing"
)

func TestNonCodeIsIgnored(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "", nil)
	assert.executeEqual(engine, "#", nil)
	assert.executeEqual(engine, "\n#", nil)
	assert.executeEqual(engine, "\t#\t\n\t#\t", nil)
	//1\n this should return nil in interactive mode
	//1\n this should return 1 in non-interactive mode
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

func TestQuantityLiterals(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "0x12ms", newDtQuantity("18", "ms"))
	assert.executeEqual(engine, "0b101ms", newDtQuantity("5", "ms"))
	assert.executeEqual(engine, "0o12ms", newDtQuantity("10", "ms"))
	assert.executeEqual(engine, "12ms", newDtQuantity("12", "ms"))
	assert.executeEqual(engine, "1.2ms", newDtQuantity("1.2", "ms"))
}

func TestValidNames(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	assert.executeEqual(engine, "a=12", newDtNumber("12"))
	assert.equal(engine.get("a"), newDtNumber("12"))
	assert.executeEqual(engine, "a_=12", newDtNumber("12"))
	assert.equal(engine.get("a_"), newDtNumber("12"))
	assert.executeEqual(engine, "a1=12", newDtNumber("12"))
	assert.equal(engine.get("a1"), newDtNumber("12"))
	assert.executeEqual(engine, "a_1=12", newDtNumber("12"))
	assert.equal(engine.get("a_1"), newDtNumber("12"))
}

func TestBasicAssigments(t *testing.T) {
	assert := newAssert(t)
	engine := newEngine()
	//spacing
	assert.executeEqual(engine, "a=12", newDtNumber("12"))
	assert.executeEqual(engine, "a \t=13", newDtNumber("13"))
	assert.executeEqual(engine, "a= \t14", newDtNumber("14"))
	assert.executeEqual(engine, "a \t= \t15", newDtNumber("15"))
	//persistance
	assert.executeEqual(engine, "a=16", newDtNumber("16"))
	assert.equal(engine.get("a"), newDtNumber("16"))
	assert.executeEqual(engine, "b=17ms", newDtQuantity("17", "ms"))
	assert.equal(engine.get("b"), newDtQuantity("17", "ms"))
}

func TestBasicExpressions(t *testing.T) {
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
