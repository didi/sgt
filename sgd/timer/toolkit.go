package timer

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"time"

	"sgt/pkg/file"
	"sgt/pkg/sys"
	"sgt/sgd/config"
)

func agentUninstall(name string) error {
	ctlfile := path.Join(config.AgsDir, name, "control")
	if !file.IsExist(ctlfile) {
		return sys.CmdRun("/bin/bash", "-c", "rm -rf "+path.Join(config.AgsDir, name))
	}

	output, err, istimeout := sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "stop")
	if istimeout {
		return fmt.Errorf("cannot stop %s: timeout", name)
	}

	if err != nil {
		return fmt.Errorf("cannot stop %s, error: %v, output: %s", name, err, output)
	}

	output, err, istimeout = sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "uninstall")
	if istimeout {
		return fmt.Errorf("cannot uninstall %s: timeout", name)
	}

	if err != nil {
		return fmt.Errorf("cannot uninstall %s, error: %v, output: %s", name, err, output)
	}

	return sys.CmdRun("/bin/bash", "-c", "rm -rf "+path.Join(config.AgsDir, name))
}

func agentInstall(name, ver string) error {
	tarFile := fmt.Sprintf("%s-%s_%s_%s.tar.gz", name, ver, runtime.GOOS, runtime.GOARCH)
	md5File := fmt.Sprintf("%s-%s_%s_%s.tar.gz.md5", name, ver, runtime.GOOS, runtime.GOARCH)
	tarPath := path.Join(config.TarDir, tarFile)
	md5Path := path.Join(config.TarDir, md5File)
	tarUrl := fmt.Sprintf("%s/tarball/%s", config.Srv, tarFile)
	md5Url := fmt.Sprintf("%s/tarball/%s", config.Srv, md5File)

	if err := ensureFilesReady(tarPath, md5Path, tarUrl, md5Url); err != nil {
		return err
	}

	dir := path.Join(config.AgsDir, name)
	if err := file.EnsureDir(dir); err != nil {
		return fmt.Errorf("cannot mkdir %s: %v", dir, err)
	}

	if err := sys.CmdRun("/bin/bash", "-c", fmt.Sprintf("tar zxf %s -C %s", tarPath, dir)); err != nil {
		return err
	}

	ctlfile := path.Join(dir, "control")
	if err := sys.CmdRun("chmod", "+x", ctlfile); err != nil {
		return fmt.Errorf("cannot chmod +x %s", ctlfile)
	}

	output, err, istimeout := sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "install")
	if istimeout {
		return fmt.Errorf("cannot install %s: timeout", name)
	}

	if err != nil {
		return fmt.Errorf("cannot install %s, error: %v, output: %s", name, err, output)
	}

	output, err, istimeout = sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "start")
	if istimeout {
		return fmt.Errorf("cannot start %s: timeout", name)
	}

	if err != nil {
		return fmt.Errorf("cannot start %s, error: %v, output: %s", name, err, output)
	}

	return nil
}

func agentUpgrade(name, ver string) error {
	if err := agentUninstall(name); err != nil {
		return err
	}

	if err := agentInstall(name, ver); err != nil {
		return err
	}

	return nil
}

func agentStart(name string) {
	ctlfile := path.Join(config.AgsDir, name, "control")
	if !file.IsExist(ctlfile) {
		return
	}

	err := sys.CmdRun("chmod", "+x", ctlfile)
	if err != nil {
		if config.Dev {
			log.Printf("INF: cannot chmod +x %s, error: %v", ctlfile, err)
		}
		return
	}

	output, err, istimeout := sys.CmdRunT(time.Duration(5)*time.Second, ctlfile, "start")
	if config.Dev {
		log.Printf("INF: %s start return output:[%s], err:[%v], istimeout:[%v]", output, err, istimeout)
	}
}

func collectAgents(dirs []string) error {
	cnt := len(dirs)
	agentStats = make(map[string]AgentStat, cnt)

	for i := 0; i < cnt; i++ {
		ctlfile := path.Join(config.AgsDir, dirs[i], "control")
		if !file.IsExist(ctlfile) {
			continue
		}

		err := sys.CmdRun("chmod", "+x", ctlfile)
		if err != nil {
			// maybe I am not root
			return fmt.Errorf("cannot chmod +x %s, error: %v", ctlfile, err)
		}

		as := AgentStat{
			Name:    dirs[i],
			Pid:     0,
			Cmdline: "",
			Ver:     "",
			Stat:    config.Stopped,
		}

		pidstr, err := sys.CmdOutTrim(ctlfile, "pid")
		if err == nil {
			pid, err := strconv.Atoi(pidstr)
			if err == nil && pid > 0 {
				as.Pid = pid
			}
		}

		if as.Pid > 0 {
			cmdline, err := file.ReadString(fmt.Sprintf("/proc/%d/cmdline", as.Pid))
			if err == nil {
				as.Cmdline = cmdline
				as.Stat = config.Started
				as.Ver, _ = sys.CmdOutTrim(ctlfile, "version")
			}
		}

		agentStats[as.Name] = as
	}

	return nil
}
