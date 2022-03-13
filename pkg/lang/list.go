package ash

import (
	"container/list"
)

type listDo struct {
	list   list.List
	length int
}

func (do *listDo) append(value interface{}) {
	do.list.PushBack(value)
	do.length++
}

func (do *listDo) each(callback func(value interface{})) {
	current := do.list.Front()
	for current != nil {
		callback(current.Value)
		current = current.Next()
	}
}
