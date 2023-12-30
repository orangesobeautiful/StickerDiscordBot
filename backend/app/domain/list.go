package domain

type ListResult[T any] struct {
	total int
	items []T
}

func (l *ListResult[T]) GetTotal() int {
	return l.total
}

func (l *ListResult[T]) GetItems() []T {
	return l.items
}

func NewListResult[T any](total int, items []T) ListResult[T] {
	return ListResult[T]{total: total, items: items}
}
