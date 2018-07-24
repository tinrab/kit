package expr

type Stack struct {
	items []interface{}
}

func NewStack() *Stack {
	return &Stack{
		items: []interface{}{},
	}
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Clear() {
	s.items = []interface{}{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	item := s.items[s.Size()-1]
	s.items = s.items[:s.Size()-1]
	return item
}

func (s *Stack) Head() interface{} {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	return s.items[s.Size()-1]
}

func (s *Stack) Tail() interface{} {
	if s.Size() == 0 {
		panic("stack is empty")
	}
	return s.items[0]
}
