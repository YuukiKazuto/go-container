package container

import (
	"errors"
)

type Stack[E comparable] interface {
	Container[E]
	GetTop() (E, error)
	Pop() (E, error)
	Push(e E)
}

var ErrStackEmpty = errors.New("error: stack cannot be empty")
