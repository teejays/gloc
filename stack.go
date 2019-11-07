package main

// Stack is a stack implementation based on an underlying array storage
type Stack struct {
	arr []interface{}
}

// NewStack returns a new instance of a Stack
func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Insert(v interface{}) error {
	s.arr = append(s.arr, v)
	return nil
}

func (s *Stack) Length() int {
	return len(s.arr)
}

func (s *Stack) Pop() interface{} {
	var elem = s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return elem
}
