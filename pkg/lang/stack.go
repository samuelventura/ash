package ash

import (
	"container/list"
)

type stackDo struct {
	list  list.List
	count int
}

func (do *stackDo) push(value interface{}) {
	do.list.PushBack(value)
	do.count++
}

func (do *stackDo) pop() interface{} {
	back := do.list.Back()
	do.list.Remove(back)
	do.count--
	return back.Value
}

func (do *stackDo) depth() int {
	return do.count
}
