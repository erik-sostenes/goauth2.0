package pkg

type Set[I comparable, V any] map[I]V

func NewSet[I comparable, V any]() Set[I, V] {
	return make(map[I]V)
}

func (s *Set[I, V]) Add(k I, v V) {
	(*s)[k] = v
}

func (s *Set[I, V]) Delete(k I) {
	delete(*s, k)
}

func (s *Set[I, V]) Size() int {
	return len(*s)
}

func (s *Set[I, V]) Exist(k I) bool {
	_, exists := (*s)[k]

	return exists
}

func (s *Set[I, V]) GetByItem(k I) V {
	return (*s)[k]
}

func (s *Set[I, V]) GetAll() []V {
	slice := make([]V, 0, len(*s))

	for _, v := range *s {
		slice = append(slice, v)
	}

	return slice
}
