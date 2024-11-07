package hijack

type seq[T any] struct {
	sequence []T
	index    int
}

type Response struct {
	Data  string
	Error error
}

func (s *seq[T]) GetNext() T {
	defer func() { s.index = s.index + 1 }()
	return s.sequence[s.index]
}

func Sequence[T any](items ...T) Sequencer[T] {
	return &seq[T]{
		sequence: items,
	}
}

type Sequencer[T any] interface {
	GetNext() T
}
