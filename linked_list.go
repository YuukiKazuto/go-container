package container

import (
	"fmt"
	"strings"
	"sync"
)

type LinkedList[E comparable] struct {
	UnimplementedLinkedContainer[E]
	head LinkedNode[E]
	len  int
	rw   sync.RWMutex
}

func NewLinkedList[E comparable](es ...E) *LinkedList[E] {
	ll := &LinkedList[E]{}
	ll.init()
	for _, e := range es {
		ll.add(e)
	}
	return ll
}

func (ll *LinkedList[E]) add(e E) {
	node := &LinkedNode[E]{
		e:    e,
		prev: ll.head.prev,
		next: &ll.head,
	}
	ll.head.prev.next = node
	ll.head.prev = node
	ll.len++
}

func (ll *LinkedList[E]) getNode(i int) *LinkedNode[E] {
	mi := ll.len / 2
	if i < mi {
		node := ll.head.next
		for j := 0; j < mi; j++ {
			if j == i {
				break
			}
			node = node.next
		}
		return node
	} else {
		node := ll.head.prev
		for j := ll.len - 1; j >= mi; j-- {
			if j == i {
				break
			}
			node = node.prev
		}
		return node
	}
}

func (ll *LinkedList[E]) init() {
	ll.head.next = &ll.head
	ll.head.prev = &ll.head
	ll.len = 0
}

func (ll *LinkedList[E]) removeNode(node *LinkedNode[E]) E {
	e := node.e
	node.prev.next = node.next
	node.next.prev = node.prev
	node = nil
	ll.len--
	return e
}

func (ll *LinkedList[E]) Head() *LinkedNode[E] {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	return &ll.head
}

func (ll *LinkedList[E]) Clear() {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	ll.init()
}

func (ll *LinkedList[E]) Get(i int) (E, error) {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= ll.len {
		return e, ErrIndexGteSize
	}
	return ll.getNode(i).e, nil
}

func (ll *LinkedList[E]) IsEmpty() bool {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	return ll.head.next == &ll.head
}

func (ll *LinkedList[E]) Iterator() Iterator[E] {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	return NewLinkedIterator[E](ll)
}

func (ll *LinkedList[E]) Size() int {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	return ll.len
}

func (ll *LinkedList[E]) ToSlice() []E {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	es := make([]E, ll.len)
	node := ll.head.next
	for i := 0; i < ll.len; i++ {
		es[i] = node.e
		node = node.next
	}
	return es
}

func (ll *LinkedList[E]) Add(e E) {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	ll.add(e)
}

func (ll *LinkedList[E]) AddToIndex(i int, e E) error {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i > ll.len {
		return ErrIndexGtSize
	}
	node := ll.getNode(i)
	n := &LinkedNode[E]{
		e:    e,
		prev: node.prev,
		next: node,
	}
	node.prev.next = n
	node.prev = n
	ll.len++
	return nil
}

func (ll *LinkedList[E]) AddList(l List[E]) error {
	if l == ll {
		return ErrSelf
	}
	ll.rw.Lock()
	defer ll.rw.Unlock()
	it := l.Iterator()
	for it.HasNext() {
		node := &LinkedNode[E]{
			e:    it.Next(),
			prev: ll.head.prev,
			next: &ll.head,
		}
		ll.head.prev.next = node
		ll.head.prev = node
		ll.len++
	}
	return nil
}

func (ll *LinkedList[E]) AddListToIndex(i int, l List[E]) error {
	if l == ll {
		return ErrSelf
	}
	ll.rw.Lock()
	defer ll.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i > ll.len {
		return ErrIndexGtSize
	}
	node := ll.getNode(i)
	it := l.Iterator()
	for it.HasNext() {
		n := &LinkedNode[E]{
			e:    it.Next(),
			prev: node.prev,
			next: node,
		}
		node.prev.next = n
		node.prev = n
		ll.len++
	}
	return nil
}

func (ll *LinkedList[E]) Copy() List[E] {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	list := &LinkedList[E]{}
	list.len = ll.len
	lNode := &list.head
	llNode := ll.head.next
	for i := 0; i < ll.len; i++ {
		lNode.next = &LinkedNode[E]{
			e:    llNode.e,
			prev: lNode,
		}
		lNode = lNode.next
		llNode = llNode.next
	}
	lNode.next = &list.head
	list.head.prev = lNode
	return list
}

func (ll *LinkedList[E]) IndexOf(e E) int {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	node := ll.head.next
	for i := 0; i < ll.len; i++ {
		if e == node.e {
			return i
		}
		node = node.next
	}
	return NotFound
}

func (ll *LinkedList[E]) LastIndexOf(e E) int {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	node := ll.head.prev
	for i := ll.len; i >= 0; i-- {
		if e == node.e {
			return i
		}
		node = node.prev
	}
	return NotFound
}

func (ll *LinkedList[E]) RemoveElements(e E) bool {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	success := false
	node := ll.head.next
	for node != &ll.head {
		next := node.next
		if e == node.e {
			ll.removeNode(node)
			success = true
		}
		node = next
	}
	return success
}

func (ll *LinkedList[E]) RemoveStart() (E, error) {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	var e E
	if ll.len == 0 {
		return e, ErrListEmpty
	}
	node := ll.head.next
	return ll.removeNode(node), nil
}

func (ll *LinkedList[E]) RemoveLast() (E, error) {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	var e E
	if ll.len == 0 {
		return e, ErrListEmpty
	}
	node := ll.head.prev
	return ll.removeNode(node), nil
}

func (ll *LinkedList[E]) RemoveByIndex(i int) (E, error) {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= ll.len {
		return e, ErrIndexGteSize
	}
	node := ll.getNode(i)
	return ll.removeNode(node), nil
}

func (ll *LinkedList[E]) Set(i int, e E) error {
	ll.rw.Lock()
	defer ll.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i >= ll.len {
		return ErrIndexGteSize
	}
	node := &ll.head
	for ii := 0; ii <= i; ii++ {
		node = node.next
	}
	node.e = e
	return nil
}

func (ll *LinkedList[E]) String() string {
	ll.rw.RLock()
	defer ll.rw.RUnlock()
	str := "["
	for node := ll.head.next; node != &ll.head; node = node.next {
		str += fmt.Sprintf("%v ", node.e)
	}
	str = strings.TrimRight(str, " ") + "]"
	return str
}
