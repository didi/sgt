package timer

import (
	"math/rand"
	"time"
)

func Sleep() {
	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
}
