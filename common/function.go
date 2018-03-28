package common

import (
	"time"
	"math/rand"
)

const (
	BACKOFF = 100
)

func CasDelay() {
	<-time.After(time.Duration(rand.Int63n(BACKOFF)) * time.Millisecond)
}

func WaitMs(last_ts int64) int64 {
	t := Ts()
	for t <= last_ts {
		t = Ts()
	}
	return t
}

// get timestamp
func Ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}