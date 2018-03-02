package utils

type Set struct {
	elements map[string]bool
	// TODO: synchronize access using a mutex
}

func NewSet(items []string) *Set {
	// convert array to map
	elements := make(map[string]bool)
	for _, name := range items {
		elements[name] = true
	}

	return &Set{
		elements: elements,
	}
}

func (s *Set) Has(element string) bool {
	return s.elements[element]
}

func (s *Set) Add(element string) {
	s.elements[element] = true
}

func (s *Set) Remove(element string) {
	delete(s.elements, element)
}

func (s *Set) Len() int {
	return len(s.elements)
}
