package cron

import (
	"log"

	"sgt/pkg/file"
)

const sgd_cron_content = `SHELL=/bin/bash
PATH=/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin:~/bin
MAILTO=root

* * * * * root timeout 3 /usr/local/sgd/control start &>/dev/null
`

func Init() {
	err := file.EnsureDir("/etc/cron.d")
	if err != nil {
		log.Printf("ERR: cannot exec mkdir -p /etc/cron.d: %v", err)
		return
	}

	_, err = file.WriteString("/etc/cron.d/sgd", sgd_cron_content)
	if err != nil {
		log.Printf("ERR: cannot write sgd cron: %v", err)
		return
	}
}
