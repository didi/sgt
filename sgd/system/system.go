package system

import (
	"hash/crc32"
	"math/rand"
	"os"
	"time"

	"sgt/sgd/config"
)

func Init() {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()+os.Getppid()) + int64(crc32.ChecksumIEEE([]byte(config.UUID))))
}
