package util

import "rsc.io/nop"

func ChannelWriteNonBlocking[T any](channel chan T, value T) {
	select {
	case channel <- value:
		nop.Nop()
	default:
		nop.Nop()
	}
}
