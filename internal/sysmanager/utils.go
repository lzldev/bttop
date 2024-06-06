package sysmanager

import (
	"time"
)

func tickerIntervalMs(interval time.Duration) time.Duration {
	return (time.Millisecond * 100) * interval
}
