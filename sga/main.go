package main

import (
	"flag"
	"log"
	"time"

	"sgt/pkg/sys"
	"sgt/pkg/watcher"
)

var Dev bool

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	isdev := flag.Bool("d", false, "in devolopment mode or not")
	flag.Parse()

	Dev = *isdev
}

func main() {
	sgdbin := "/usr/local/sgd/bin/sgd"
	ctlfile := "/usr/local/sgd/control"

	w := watcher.NewWatcher(sgdbin, time.Duration(7)*time.Second)
	w.Dev = Dev
	w.Start(func() {
		output, err, istimeout := sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "start")
		if istimeout {
			log.Printf("cannot start %s, timeout", sgdbin)
			return
		}

		if err != nil {
			log.Printf("cannot start %s, error: %v, output: %s", sgdbin, err, output)
		}
	})
}
