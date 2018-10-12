package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"sgt/pkg/logger"
	"sgt/pkg/tools"
	"sgt/sgd/config"
	"sgt/sgd/cron"
	"sgt/sgd/system"
	"sgt/sgd/timer"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	server := flag.String("s", "http://127.0.0.1:8030", "server address")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	isdev := flag.Bool("d", false, "in devolopment mode or not")
	puuid := flag.Bool("u", false, "print vm uuid")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleServer(*server)
	handleIsdev(*isdev)
	handlePrint(*puuid)
}

func main() {
	cron.Init()
	config.Init()

	logDir := path.Join(config.Cwd, "log")
	lb, err := logger.NewFileBackend(logDir)
	if err != nil {
		log.Fatalln("FAT: cannot init logger:", err)
	}

	logger.SetLogging("ERROR", lb)
	lb.Rotate(config.LogFileNum, config.LogFileSize)

	defer func() {
		logger.Close()
	}()

	tools.OnInterrupt(func() {
		logger.Close()
		os.Exit(0)
	})

	system.Init()
	timer.Init()

	timer.Sleep()

	go timer.Heartbeat()
	go timer.SgaWatch()

	select {}
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(config.Ver)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleServer(addr string) {
	config.Srv = addr
}

func handleIsdev(isdev bool) {
	config.Dev = isdev
}

func handlePrint(pu bool) {
	if pu {
		uuid, err := config.GetUUID()
		if err != nil {
			os.Exit(1)
		}

		fmt.Println(uuid)
		os.Exit(0)
	}
}
