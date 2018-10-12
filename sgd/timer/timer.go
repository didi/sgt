package timer

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"sgt/pkg/file"
	"sgt/sgd/config"
)

var httpcli http.Client

func Init() {
	initHttpcli()
	initDirs()
	initAgents()
}

func initHttpcli() {
	httpcli = http.Client{
		Timeout: time.Second * 5,
	}
}

func initDirs() {
	if err := file.EnsureDir(config.AgsDir); err != nil {
		log.Fatalf("FAT: cannot mkdir %s, error: %v", config.AgsDir, err)
	}

	if err := file.EnsureDir(config.LogDir); err != nil {
		log.Fatalf("FAT: cannot mkdir %s, error: %v", config.LogDir, err)
	}

	if err := file.EnsureDir(config.TarDir); err != nil {
		log.Fatalf("FAT: cannot mkdir %s, error: %v", config.TarDir, err)
	}
}

func initAgents() {
	agentDirs, err := file.DirsUnder(config.AgsDir)
	if err != nil {
		log.Fatalf("FAT: cannot list %s, error: %v", config.AgsDir, err)
	}

	sort.Strings(agentDirs)
	lastAgentDirs = strings.Join(agentDirs, "")

	err = collectAgents(agentDirs)
	if err != nil {
		log.Fatalf("FAT: %v", err)
	}
}
