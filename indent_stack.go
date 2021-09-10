package tlps

import "github.com/goropikari/tlps/collections/stack"

type IndentStack struct {
	Stack *stack.Stack
}

func NewIndentStack() *IndentStack {
	return &IndentStack{stack.NewStack()}
}

// Push adds an item in stack
func (s *IndentStack) Push(x int) {
	s.Stack.Push(x)
}

// Pop pops an item from stack
func (s *IndentStack) Pop() int {
	return s.Stack.Pop().(int)
}

// Peek returns top item in stack, and don't modity the stack.
func (s *IndentStack) Peek() int {
	return s.Stack.Peek().(int)
}

// IsEmpty checks that stack is empty
func (s *IndentStack) IsEmpty() bool {
	return s.Stack.IsEmpty()
}

// Size returns stack size
func (s *IndentStack) Size() int {
	return s.Stack.Size()
}
