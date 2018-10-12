package timer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"sgt/pkg/file"
)

func destruct() {
	err := file.Unlink("/etc/cron.d/sgd")
	if err != nil {
		log.Println("cannot del /etc/cron.d/sgd:", err)
	}

	fs, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Println("cannot read /proc:", err)
		os.Exit(1)
	}

	sz := len(fs)
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			continue
		}

		name := fs[i].Name()
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		exe := fmt.Sprintf("/proc/%d/exe", pid)
		if !file.IsExist(exe) {
			continue
		}

		target, err := os.Readlink(exe)
		if err == nil && strings.Contains(target, "/usr/local/sgd/bin/sga") {
			proc, err := os.FindProcess(pid)
			if err != nil {
				continue
			}

			err = proc.Kill()
			if err != nil {
				log.Printf("ERR: cannot kill process[pid:%d]: %v", pid, err)
			}
		}
	}

	os.Exit(1)
}
