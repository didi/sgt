package timer

import (
	"log"
	"runtime"

	"sgt/sgd/config"
)

const (
	max_mem uint64 = 30 * 1024 * 1024
)

func checkMem() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	if config.Dev {
		log.Printf("INF: mem use: %dMB", m.HeapSys/1024/1024)
	}

	if m.HeapSys > max_mem {
		log.Fatalf("FAT: mem use: %dMB, overload", m.HeapSys/1024/1024)
	}
}
