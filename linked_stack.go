package container

import (
	"fmt"
	"strings"
	"sync"
)

type LinkedStack[E comparable] struct {
	UnimplementedLinkedContainer[E]
	head LinkedNode[E]
	len  int
	rw   sync.RWMutex
}

func NewLinkedStack[E comparable](es ...E) *LinkedStack[E] {
	ls := &LinkedStack[E]{}
	ls.init()
	for _, e := range es {
		ls.push(e)
	}
	return ls
}

func (ls *LinkedStack[E]) getNode(i int) *LinkedNode[E] {
	mi := ls.len / 2
	if i < mi {
		node := ls.head.next
		for j := 0; j < mi; j++ {
			if j == i {
				break
			}
			node = node.next
		}
		return node
	} else {
		node := ls.head.prev
		for j := ls.len - 1; j >= mi; j-- {
			if j == i {
				break
			}
			node = node.prev
		}
		return node
	}
}

func (ls *LinkedStack[E]) init() {
	ls.head.prev = &ls.head
	ls.head.next = &ls.head
	ls.len = 0
}

func (ls *LinkedStack[E]) push(e E) {
	node := &LinkedNode[E]{
		e:    e,
		prev: ls.head.prev,
		next: &ls.head,
	}
	ls.head.prev.next = node
	ls.head.prev = node
	ls.len++
}

func (ls *LinkedStack[E]) Head() *LinkedNode[E] {
	return &ls.head
}

func (ls *LinkedStack[E]) Clear() {
	ls.rw.Lock()
	defer ls.rw.Unlock()
	ls.init()
}

func (ls *LinkedStack[E]) Get(i int) (E, error) {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= ls.len {
		return e, ErrIndexGteSize
	}
	e = ls.getNode(i).e
	return e, nil
}

func (ls *LinkedStack[E]) IsEmpty() bool {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	return ls.head.next == &ls.head
}

func (ls *LinkedStack[E]) Iterator() Iterator[E] {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	return NewLinkedIterator[E](ls)
}

func (ls *LinkedStack[E]) Size() int {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	return ls.len
}

func (ls *LinkedStack[E]) ToSlice() []E {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	es := make([]E, ls.len)
	node := ls.head.next
	for i := 0; i < ls.len; i++ {
		es[i] = node.e
		node = node.next
	}
	return es
}

func (ls *LinkedStack[E]) GetTop() (E, error) {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	var e E
	top := ls.head.prev
	if top == &ls.head {
		return e, ErrStackEmpty
	}
	return top.e, nil
}

func (ls *LinkedStack[E]) Pop() (E, error) {
	ls.rw.Lock()
	defer ls.rw.Unlock()
	var e E
	top := ls.head.prev
	if top == &ls.head {
		return e, ErrStackEmpty
	}
	e = top.e
	top.prev.next = &ls.head
	ls.head.prev = top.prev
	top = nil
	return e, nil
}

func (ls *LinkedStack[E]) Push(e E) {
	ls.rw.Lock()
	defer ls.rw.Unlock()
	ls.push(e)
}

func (ls *LinkedStack[E]) String() string {
	ls.rw.RLock()
	defer ls.rw.RUnlock()
	str := "["
	for node := ls.head.next; node != &ls.head; node = node.next {
		str += fmt.Sprintf("%v ", node.e)
	}
	str = strings.TrimRight(str, " ") + "]"
	return str
}
