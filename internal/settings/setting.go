package settings

func New[T any](name string, defaultv, value T, desc string, states ...T) *Setting[T] {
	return &Setting[T]{
		name:        name,
		defaultv:    defaultv,
		value:       value,
		description: desc,
	}
}

type Setting[T any] struct {
	name        string
	value       T
	states      []T
	defaultv    T
	description string
}

func (self *Setting[T]) Name() string {
	return self.name
}

func (self *Setting[T]) Get() T {
	return self.value
}

func (self *Setting[T]) Set(v T) {
	self.value = v
}

var (
	Int    = New[int]
	Bool   = New[bool]
	String = New[string]
)
