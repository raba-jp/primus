package repl

import "io"

type State struct {
	Continuation bool
	Inputs       []string
	next         int
}

func NewState() *State {
	return &State{
		Continuation: false,
		Inputs:       make([]string, 0, 1),
		next:         0,
	}
}

func (s *State) Readline() ([]byte, error) {
	if s.next >= len(s.Inputs) {
		s.next = 0
		s.Continuation = true
		return nil, io.EOF
	}
	val := s.Inputs[s.next]
	s.next++
	return []byte(val), nil
}

func (s *State) AppendInput(val string) {
	s.Inputs = append(s.Inputs, val+"\n")
}

func (s *State) Reset() {
	s.Continuation = false
	s.Inputs = make([]string, 1)
	s.next = 0
}
