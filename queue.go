package container

import "errors"

type Queue[E comparable] interface {
	Container[E]
	En(e E)
	De() (E, error)
	GetFront() (E, error)
	GetRear() (E, error)
}

var ErrQueueEmpty = errors.New("error: queue cannot be empty")
