package common

import (
	"testing"
)

type testStruct struct {
	A string
	B string
}

func NewtestStruct() testStruct {
	return testStruct{
		A: "-",
		B: "-",
	}
}

func TestSinglyLinkedList(t *testing.T) {
	l := NewSinglyLinkedList[int](nil, nil, 5)

	if l.RemoveNodeAtHead() != nil {
		t.Errorf("Non-nil response for getting head node from empty list.\n")
	}

	if l.Len != 0 {
		t.Errorf("Expected zero length, got %d\n", l.Len)
	}

	l.AddNodeAtTail(3)
	l.AddNodeAtTail(2)
	l.AddNodeAtTail(1)

	if l.Len != 3 {
		t.Errorf("Expected length of 3, got %d\n", l.Len)
	}

	if ret := l.RemoveNodeAtHead(); ret == nil {
		t.Errorf("Expected node (3) got nil value %v\n", ret)
	} else if ret.Value != 3 {
		t.Errorf("Expected value 3 got %v\n", ret)
	}

	if ret := l.RemoveNodeAtHead(); ret == nil {
		t.Errorf("Expected node (2) got a nil value %v\n", ret)
	} else if ret.Value != 2 {
		t.Errorf("Expected value 2 got %v\n", ret)
	}

	if ret := l.RemoveNodeAtHead(); ret == nil {
		t.Errorf("Expected node (1) got a nil value %v\n", ret)
	} else if ret.Value != 1 {
		t.Errorf("Expected value 1 got %v\n", ret)
	}

	if l.Len != 0 {
		t.Errorf("Expected zero length, got %d\n", l.Len)
	}

	for i := 0; i <= 5; i++ {
		l.AddNodeAtTail(i)
	}

	if ret := l.RemoveNodeAtHead(); ret.Value != 1 {
		t.Errorf("Expected value 1 got %v\n", ret)
	}

	l.AddNodeAtTail(6)

	if ret := l.RemoveNodeAtHead(); ret.Value != 2 {
		t.Errorf("Expected value 2 got %v\n", ret)
	}

	ll := NewSinglyLinkedList[testStruct](nil, nil, 5)
	ll.AddNodeAtTail(NewtestStruct())
	ll.AddNodeAtTail(NewtestStruct())

	if ret := ll.RemoveNodeAtHead(); ret.Value.A != "-" && ret.Value.B != "-" {
		t.Errorf("Expected - and -, got %s and %s\n", ret.Value.A, ret.Value.B)
	}
}
