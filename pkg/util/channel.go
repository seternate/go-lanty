package util

func ChannelWriteNonBlocking[T any](channel chan T, value T) {
	select {
	case channel <- value:
		return
	default:
		return
	}
}
