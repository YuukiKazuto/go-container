package container

// LinkedNode 此结构体定义双向链表节点结构
// 链表，链式栈，链式队列的节点皆可使用此结构
type LinkedNode[E comparable] struct {
	e          E
	prev, next *LinkedNode[E]
}

func (node *LinkedNode[E]) Value() E {
	return node.e
}

func (node *LinkedNode[E]) SetValue(e E) {
	node.e = e
}

func (node *LinkedNode[E]) Prev() *LinkedNode[E] {
	return node.prev
}

func (node *LinkedNode[E]) SetPrev(prev *LinkedNode[E]) {
	node.prev = prev
}

func (node *LinkedNode[E]) Next() *LinkedNode[E] {
	return node.next
}

func (node *LinkedNode[E]) SetNext(next *LinkedNode[E]) {
	node.next = next
}
