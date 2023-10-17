package container

import (
	"errors"
)

type List[E comparable] interface {
	Container[E]
	Add(e E)
	AddToIndex(i int, e E) error
	AddList(l List[E]) error
	AddListToIndex(i int, l List[E]) error
	Copy() List[E]
	IndexOf(e E) int
	LastIndexOf(e E) int
	RemoveElements(e E) bool
	RemoveStart() (E, error)
	RemoveLast() (E, error)
	RemoveByIndex(i int) (E, error)
	Set(i int, e E) error
}

const NotFound = -1

var (
	ErrListEmpty = errors.New("error: list cannot be empty")
	ErrSelf      = errors.New("error: param l cannot be the caller itself")
)
