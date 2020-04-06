package stack

import (
	"container/list"
	"sync"
)

type Stack struct {
	lock  *sync.Mutex
	stack *list.List
}

func Init() *Stack {
	return &Stack{
		&sync.Mutex{},
		list.New(),
	}
}

func (s *Stack) Push(x interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stack.PushFront(x)
}

func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.stack.Len() > 0 {
		ele := s.stack.Front()
		s.stack.Remove(ele)
		return ele.Value
	}
	return nil
}

func (s *Stack) Front() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.stack.Len() > 0 {
		return s.stack.Front().Value
	}
	return nil
}

func (s *Stack) Size() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.stack.Len()
}

func (s *Stack) Empty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.stack.Len() == 0
}
