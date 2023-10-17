package container

import "errors"

type Iterator[E any] interface {
	HasNext() bool
	Next() E
}

type Container[E any] interface {
	Clear()
	Get(i int) (E, error)
	IsEmpty() bool
	Iterator() Iterator[E]
	Size() int
	ToSlice() []E
}

var (
	ErrIndexGtSize  = errors.New("error: param i cannot be greater than the container size")
	ErrIndexGteSize = errors.New("error: param i cannot be greater than or equal to the container size")
	ErrIndexLtZero  = errors.New("error: param i cannot be less than 0")
)

// SqContainer is sequence container
// All implementations must embed UnimplementedSqContainer
type SqContainer[E comparable] interface {
	Get(i int) (E, error)
	Size() int
	mustEmbedUnimplementedSqContainer()
}

type UnimplementedSqContainer[E comparable] struct{}

func (UnimplementedSqContainer[E]) Get(i int) (E, error) {
	var e E
	return e, errors.New("method Get not implemented")
}

func (UnimplementedSqContainer[E]) Size() int {
	return 0
}

func (UnimplementedSqContainer[E]) mustEmbedUnimplementedSqContainer() {}

type SqIterator[E comparable] struct {
	sc    SqContainer[E]
	index int
}

func NewSqIterator[E comparable](sc SqContainer[E]) *SqIterator[E] {
	return &SqIterator[E]{sc: sc}
}

func (it *SqIterator[E]) HasNext() bool {
	return it.index < it.sc.Size()
}

func (it *SqIterator[E]) Next() E {
	e, _ := it.sc.Get(it.index)
	it.index++
	return e
}

// LinkedContainer is linked container
// All implementations must embed UnimplementedLinkedContainer
type LinkedContainer[E comparable] interface {
	Head() *LinkedNode[E]
	mustEmbedUnimplementedLinkedContainer()
}

type UnimplementedLinkedContainer[E comparable] struct{}

func (UnimplementedLinkedContainer[E]) Head() *LinkedNode[E] {
	return nil
}

func (UnimplementedLinkedContainer[E]) mustEmbedUnimplementedLinkedContainer() {}

type LinkedIterator[E comparable] struct {
	lc      LinkedContainer[E]
	curNode *LinkedNode[E]
}

func NewLinkedIterator[E comparable](lc LinkedContainer[E]) *LinkedIterator[E] {
	return &LinkedIterator[E]{
		lc:      lc,
		curNode: lc.Head(),
	}
}

func (l *LinkedIterator[E]) HasNext() bool {
	return l.curNode.Next() != l.lc.Head()
}

func (l *LinkedIterator[E]) Next() E {
	next := l.curNode.Next()
	if next != l.lc.Head() {
		l.curNode = next
	}
	return next.Value()
}
