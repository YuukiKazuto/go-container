package container

import (
	"fmt"
	"sync"
)

type SliceList[E comparable] struct {
	UnimplementedSqContainer[E]
	elems []E
	rw    sync.RWMutex
}

func NewSliceList[E comparable](es ...E) *SliceList[E] {
	return &SliceList[E]{
		elems: es,
	}
}

func (sl *SliceList[E]) Clear() {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	sl.elems = nil
}

func (sl *SliceList[E]) Get(i int) (E, error) {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= len(sl.elems) {
		return e, ErrIndexGteSize
	}
	return sl.elems[i], nil
}

func (sl *SliceList[E]) IsEmpty() bool {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	return len(sl.elems) == 0
}

func (sl *SliceList[E]) Iterator() Iterator[E] {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	return NewSqIterator[E](sl)
}

func (sl *SliceList[E]) Size() int {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	return len(sl.elems)
}

func (sl *SliceList[E]) ToSlice() []E {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	es := make([]E, len(sl.elems))
	copy(es, sl.elems)
	return es
}

func (sl *SliceList[E]) Add(e E) {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	sl.elems = append(sl.elems, e)
}

func (sl *SliceList[E]) AddToIndex(i int, e E) error {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i > len(sl.elems) {
		return ErrIndexGtSize
	}
	sl.elems = append(sl.elems, e)
	for j := len(sl.elems); j > i; j-- {
		sl.elems[j] = sl.elems[j-1]
	}
	sl.elems[i] = e
	return nil
}

func (sl *SliceList[E]) AddList(l List[E]) error {
	if l == sl {
		return ErrSelf
	}
	sl.rw.Lock()
	defer sl.rw.Unlock()
	if a, ok := l.(*SliceList[E]); ok {
		sl.elems = append(sl.elems, a.elems...)
	} else {
		iterator := l.Iterator()
		for iterator.HasNext() {
			sl.elems = append(sl.elems, iterator.Next())
		}
	}
	return nil
}

func (sl *SliceList[E]) AddListToIndex(i int, l List[E]) error {
	if l == sl {
		return ErrSelf
	}
	sl.rw.Lock()
	defer sl.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i > len(sl.elems) {
		return ErrIndexGteSize
	}
	if a, ok := l.(*SliceList[E]); ok {
		sl.elems = append(sl.elems[:i], append(a.elems, sl.elems[i:]...)...)
	} else {
		a2 := make([]E, len(sl.elems[i:]))
		copy(a2, sl.elems[i:])
		iterator := l.Iterator()
		for j := i; iterator.HasNext(); j++ {
			sl.elems = append(sl.elems[:j], iterator.Next())
		}
		sl.elems = append(sl.elems, a2...)
	}
	return nil
}

func (sl *SliceList[E]) Copy() List[E] {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	list := &SliceList[E]{}
	list.elems = make([]E, len(sl.elems))
	copy(list.elems, sl.elems)
	return list
}

func (sl *SliceList[E]) IndexOf(e E) int {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	for i, v := range sl.elems {
		if v == e {
			return i
		}
	}
	return NotFound
}

func (sl *SliceList[E]) LastIndexOf(e E) int {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	for i := len(sl.elems); i >= 0; i-- {
		if e == sl.elems[i] {
			return i
		}
	}
	return NotFound
}

func (sl *SliceList[E]) RemoveElements(e E) bool {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	success := false
	n := len(sl.elems)
	for i := 0; i < n; {
		if e == sl.elems[i] {
			sl.elems = append(sl.elems[:i], sl.elems[i+1:]...)
			n--
			success = true
			continue
		}
		i++
	}
	return success
}

func (sl *SliceList[E]) RemoveStart() (E, error) {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	var e E
	if len(sl.elems) == 0 {
		return e, ErrListEmpty
	}
	i := len(sl.elems) - 1
	e = sl.elems[i]
	sl.elems = sl.elems[:i]
	return e, nil
}

func (sl *SliceList[E]) RemoveLast() (E, error) {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	var e E
	if len(sl.elems) == 0 {
		return e, ErrListEmpty
	}
	i := len(sl.elems) - 1
	e = sl.elems[i]
	sl.elems = sl.elems[:i]
	return e, nil
}

func (sl *SliceList[E]) RemoveByIndex(i int) (E, error) {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	var e E
	if i < 0 {
		return e, ErrIndexLtZero
	}
	if i >= len(sl.elems) {
		return e, ErrIndexGteSize
	}
	e = sl.elems[i]
	sl.elems = append(sl.elems[:i], sl.elems[i+1:]...)
	return e, nil
}

func (sl *SliceList[E]) Set(i int, e E) error {
	sl.rw.Lock()
	defer sl.rw.Unlock()
	if i < 0 {
		return ErrIndexLtZero
	}
	if i >= len(sl.elems) {
		return ErrIndexGteSize
	}
	sl.elems[i] = e
	return nil
}

func (sl *SliceList[E]) String() string {
	sl.rw.RLock()
	defer sl.rw.RUnlock()
	return fmt.Sprint(sl.elems)
}
