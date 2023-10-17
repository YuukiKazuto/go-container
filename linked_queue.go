package container

import (
	"fmt"
	"strings"
	"sync"
)

type LinkedQueue[E comparable] struct {
	UnimplementedLinkedContainer[E]
	head LinkedNode[E]
	len  int
	rw   sync.RWMutex
}

func NewLinkedQueue[E comparable](es ...E) *LinkedQueue[E] {
	ls := &LinkedQueue[E]{}
	ls.init()
	for _, e := range es {
		ls.en(e)
	}
	return ls
}

func (lq *LinkedQueue[E]) getNode(i int) *LinkedNode[E] {
	mi := lq.len / 2
	if i < mi {
		node := lq.head.next
		for j := 0; j < mi; j++ {
			if j == i {
				break
			}
			node = node.next
		}
		return node
	} else {
		node := lq.head.prev
		for j := lq.len - 1; j >= mi; j-- {
			if j == i {
				break
			}
			node = node.prev
		}
		return node
	}
}

func (lq *LinkedQueue[E]) init() {
	lq.head.prev = &lq.head
	lq.head.next = &lq.head
	lq.len = 0
}

func (lq *LinkedQueue[E]) en(e E) {
	node := &LinkedNode[E]{
		e:    e,
		prev: lq.head.prev,
		next: &lq.head,
	}
	lq.head.prev.next = node
	lq.head.prev = node
	lq.len++
}

func (lq *LinkedQueue[E]) Head() *LinkedNode[E] {
	return &lq.head
}

func (lq *LinkedQueue[E]) Clear() {
	lq.rw.Lock()
	defer lq.rw.Unlock()
	lq.init()
}

func (lq *LinkedQueue[E]) Get(i int) (E, error) {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= lq.len {
		return e, ErrIndexGteSize
	}
	e = lq.getNode(i).e
	return e, nil
}

func (lq *LinkedQueue[E]) IsEmpty() bool {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	return lq.head.next == &lq.head
}

func (lq *LinkedQueue[E]) Iterator() Iterator[E] {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	return NewLinkedIterator[E](lq)
}

func (lq *LinkedQueue[E]) Size() int {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	return lq.len
}

func (lq *LinkedQueue[E]) ToSlice() []E {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	es := make([]E, lq.len)
	node := lq.head.next
	for i := 0; i < lq.len; i++ {
		es[i] = node.e
		node = node.next
	}
	return es
}

func (lq *LinkedQueue[E]) En(e E) {
	lq.rw.Lock()
	defer lq.rw.Unlock()
	lq.en(e)
}

func (lq *LinkedQueue[E]) De() (E, error) {
	lq.rw.Lock()
	defer lq.rw.Unlock()
	var e E
	front := lq.head.next
	if front == &lq.head {
		return e, ErrStackEmpty
	}
	e = front.e
	front.next.prev = &lq.head
	lq.head.next = front.next
	front = nil
	return e, nil
}

func (lq *LinkedQueue[E]) GetFront() (E, error) {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	var e E
	front := lq.head.next
	if front == &lq.head {
		return e, ErrQueueEmpty
	}
	return front.e, nil
}

func (lq *LinkedQueue[E]) GetRear() (E, error) {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	var e E
	rear := lq.head.prev
	if rear == &lq.head {
		return e, ErrQueueEmpty
	}
	return rear.e, nil
}

func (lq *LinkedQueue[E]) String() string {
	lq.rw.RLock()
	defer lq.rw.RUnlock()
	str := "["
	for node := lq.head.next; node != &lq.head; node = node.next {
		str += fmt.Sprintf("%v ", node.e)
	}
	str = strings.TrimRight(str, " ") + "]"
	return str
}
