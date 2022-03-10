package ash

import "testing"

func TestBasicTree(t *testing.T) {
	script := `
	tree global:
		count = 0
		increment:
			count++
		add amount:
			count += amount
	`
	engine := Compile(script)
	engine.assert_equal("global.count", 0)
	engine.execute("global.increment")
	engine.assert_equal("global.count", 1)
	engine.execute("global.add 2")
	engine.assert_equal("global.count", 3)
}
