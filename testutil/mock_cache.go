package testutil

type MockCache[T any] struct {
	SetFunc    func(key string, val T)
	GetFunc    func(key string) (T, bool)
	DeleteFunc func(key string)
}

func (m *MockCache[T]) Set(key string, val T) {
	if m.SetFunc != nil {
		m.SetFunc(key, val)
	}
}

func (m *MockCache[T]) Get(key string) (T, bool) {
	if m.GetFunc == nil {
		var zero T
		return zero, false
	}
	return m.GetFunc(key)
}

func (m *MockCache[T]) Delete(key string) {
	if m.DeleteFunc != nil {
		m.DeleteFunc(key)
	}
}
