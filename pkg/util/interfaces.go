package util

type Publisher[T any] interface {
	Subscribe(chan T)
	Unsubscribe(chan T)
}
