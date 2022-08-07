package common

/*
This package contains a custom (partial) implementation of a singly linked list which
accepts a generic type for the value/data it holds. The list also has len and capacity
fields which helps in it being used as a queue.
*/

// Node for a singly linked list
type SinglyNode[T any] struct {
	Next  *SinglyNode[T]
	Value T
}

// NewSinglyNode returns a pointer to a node
func NewSinglyNode[T any](next *SinglyNode[T], val T) *SinglyNode[T] {
	return &SinglyNode[T]{
		Next:  next,
		Value: val,
	}
}

//
type SinglyLinkedList[T any] struct {
	Head, Tail    *SinglyNode[T]
	Capacity, Len int
}

func NewSinglyLinkedList[T any](head, tail *SinglyNode[T], cap int) *SinglyLinkedList[T] {
	if cap == 0 {
		return nil
	}
	return &SinglyLinkedList[T]{
		Head:     head,
		Tail:     tail,
		Capacity: cap,
		Len:      0,
	}
}

func (l *SinglyLinkedList[T]) RemoveNodeAtHead() *SinglyNode[T] {
	// Empty list
	if l.Head == nil {
		return nil
	}

	ret := l.Head
	l.Head = l.Head.Next

	// Single element list
	if ret == l.Tail {
		l.Tail = nil
	}

	l.Len -= 1

	return ret
}

func (l *SinglyLinkedList[T]) AddNodeAtTail(val T) {
	n := NewSinglyNode(nil, val)

	// Empty list
	if l.Tail != nil {
		l.Tail.Next = n
	}

	l.Tail = n
	l.Len += 1

	// Empty linked list check
	if l.Head == nil {
		l.Head = l.Tail
	}

	if l.Len > l.Capacity {
		l.RemoveNodeAtHead()
	}
}
