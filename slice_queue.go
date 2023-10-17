package container

import (
	"fmt"
	"sync"
)

type SliceQueue[E comparable] struct {
	UnimplementedSqContainer[E]
	elems []E
	rw    sync.RWMutex
}

func NewSliceQueue[E comparable](es ...E) *SliceQueue[E] {
	return &SliceQueue[E]{elems: es}
}

func (sq *SliceQueue[E]) Clear() {
	sq.rw.Lock()
	defer sq.rw.Unlock()
	sq.elems = nil
}

func (sq *SliceQueue[E]) Get(i int) (E, error) {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= len(sq.elems) {
		return e, ErrIndexGteSize
	}
	return sq.elems[i], nil
}

func (sq *SliceQueue[E]) IsEmpty() bool {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	return len(sq.elems) == 0
}

func (sq *SliceQueue[E]) Iterator() Iterator[E] {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	return NewSqIterator[E](sq)
}

func (sq *SliceQueue[E]) Size() int {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	return len(sq.elems)
}

func (sq *SliceQueue[E]) ToSlice() []E {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	es := make([]E, len(sq.elems))
	copy(es, sq.elems)
	return es
}

func (sq *SliceQueue[E]) En(e E) {
	sq.rw.Lock()
	defer sq.rw.Unlock()
	sq.elems = append(sq.elems, e)
}

func (sq *SliceQueue[E]) De() (E, error) {
	sq.rw.Lock()
	defer sq.rw.Unlock()
	var e E
	if len(sq.elems) == 0 {
		return e, ErrQueueEmpty
	}
	e = sq.elems[0]
	sq.elems = sq.elems[1:]
	return e, nil
}

func (sq *SliceQueue[E]) GetFront() (E, error) {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	var e E
	if len(sq.elems) == 0 {
		return e, ErrQueueEmpty
	}
	e = sq.elems[0]
	return e, nil
}

func (sq *SliceQueue[E]) GetRear() (E, error) {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	var e E
	n := len(sq.elems)
	if n == 0 {
		return e, ErrQueueEmpty
	}
	e = sq.elems[n-1]
	return e, nil
}

func (sq *SliceQueue[E]) String() string {
	sq.rw.RLock()
	defer sq.rw.RUnlock()
	return fmt.Sprint(sq.elems)
}
