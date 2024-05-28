package stack

type Stack[T any] struct {
  arr []T
  size int
}

func NewStack[T any]() *Stack[T] {
  return &Stack[T]{}
}

func (s *Stack[T]) Empty() bool {
	return s.size == 0
}

func (s *Stack[T]) Top() T {
	return s.arr[s.size - 1]
}

func (s *Stack[T]) Push(val T) {
  if cap(s.arr) == s.size {
    s.arr = append(s.arr, val)
  } else {
    s.arr[s.size] = val
  }
  s.size++
}

func (s *Stack[T]) Pop() T {
  s.size--
  return s.arr[s.size]
}

