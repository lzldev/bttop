package sysmanager

import (
	"iter"
	"time"
)

func tickerIntervalMs(interval time.Duration) time.Duration {
	return (time.Millisecond * 100) * interval
}

func Map[TIn any, TOut any](i iter.Seq[TIn], mapFn func(TIn) TOut) iter.Seq[TOut] {
	return func(yield func(TOut) bool) {
		next, stop := iter.Pull(i)
		defer stop()

		for {
			v1, ok1 := next()
			if !ok1 {
				return
			}
			if !yield(mapFn(v1)) {
				return
			}
		}
	}
}
