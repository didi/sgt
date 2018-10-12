package timer

import (
	"fmt"
	"time"

	"sgt/pkg/logger"
	"sgt/pkg/sys"
	"sgt/pkg/watcher"
	"sgt/sgd/config"
)

func SgaWatch() {
	sgabin := "/usr/local/sgd/bin/sga"

	w := watcher.NewWatcher(sgabin, time.Duration(7)*time.Second)
	w.Dev = config.Dev
	w.Start(func() {
		output, err, istimeout := sys.CmdRunT(time.Duration(5)*time.Second, "/bin/bash", "-c", fmt.Sprintf("nohup %s &>/dev/null &", sgabin))
		if istimeout {
			logger.Errorf("cannot start %s, timeout", sgabin)
			return
		}

		if err != nil {
			logger.Errorf("cannot start %s, error: %v, output: %s", sgabin, err, output)
		}
	})
}
