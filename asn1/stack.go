package asn1

type stack struct {
	data []interface{}	
}

func (s *stack)Push(v interface{}) {
	s.data = append(s.data, v)
}


func (s *stack)Pop() interface{} {
	var v interface{}
	l := len(s.data)
	if l > 0 {
		l--
		v = s.data[l]
		s.data = s.data[:l]
		return v
	}
	return nil
}

func NewStack() *stack {
	return &stack{}
}