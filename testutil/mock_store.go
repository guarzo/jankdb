package testutil

// MockStore is a partial stand-in for jankdb.Store[T],
// if you prefer an interface-based approach.
// But typically you'd just use a real Store[T] with a MockFileSystem in your tests.
type MockStore[T any] struct {
	LoadFunc func() error
	SaveFunc func() error
	GetFunc  func() T
	SetFunc  func(T)
}

func (m *MockStore[T]) Load() error {
	if m.LoadFunc != nil {
		return m.LoadFunc()
	}
	return nil
}
func (m *MockStore[T]) Save() error {
	if m.SaveFunc != nil {
		return m.SaveFunc()
	}
	return nil
}
func (m *MockStore[T]) Get() T {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	var zero T
	return zero
}
func (m *MockStore[T]) Set(val T) {
	if m.SetFunc != nil {
		m.SetFunc(val)
	}
}
