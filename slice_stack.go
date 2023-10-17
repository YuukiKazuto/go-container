package container

import (
	"fmt"
	"sync"
)

type SliceStack[E comparable] struct {
	UnimplementedSqContainer[E]
	elems []E
	rw    sync.RWMutex
}

func NewSliceStack[E comparable](es ...E) *SliceStack[E] {
	return &SliceStack[E]{
		elems: es,
	}
}

func (ss *SliceStack[E]) Clear() {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	ss.elems = nil
}

func (ss *SliceStack[E]) Get(i int) (E, error) {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= len(ss.elems) {
		return e, ErrIndexGteSize
	}
	return ss.elems[i], nil
}

func (ss *SliceStack[E]) IsEmpty() bool {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	return len(ss.elems) == 0
}

func (ss *SliceStack[E]) Iterator() Iterator[E] {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	return NewSqIterator[E](ss)
}

func (ss *SliceStack[E]) Size() int {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	return len(ss.elems)
}

func (ss *SliceStack[E]) ToSlice() []E {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	es := make([]E, len(ss.elems))
	copy(es, ss.elems)
	return es
}

func (ss *SliceStack[E]) GetTop() (E, error) {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	var e E
	n := len(ss.elems)
	if n == 0 {
		return e, ErrStackEmpty
	}
	e = ss.elems[n-1]
	return e, nil
}

func (ss *SliceStack[E]) Pop() (E, error) {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	var e E
	n := len(ss.elems)
	if n == 0 {
		return e, ErrStackEmpty
	}
	e = ss.elems[n-1]
	ss.elems = ss.elems[:n-1]
	return e, nil
}

func (ss *SliceStack[E]) Push(e E) {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	ss.elems = append(ss.elems, e)
}

func (ss *SliceStack[E]) String() string {
	return fmt.Sprint(ss.elems)
}
