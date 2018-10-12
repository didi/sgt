package timer

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"sgt/pkg/file"
	"sgt/pkg/logger"
	"sgt/pkg/sys"
	"sgt/sgd/config"
)

func handleHeartbeatResponse(res HeartbeatResponse) {
	if res.Error != "" {
		logger.Error("heartbeat server return error: ", res.Error)
		return
	}

	if res.SgdVer == "-1" {
		// start self destruct program
		destruct()
		return
	}

	if res.SgdVer != config.Ver {
		upgradeSgd(res.SgdVer)
		return
	}

	digest := agentVersDigest(res)
	if lastAgentVers == digest {
		return
	}

	if err := handleAgVers(res.Ags); err != nil {
		logger.Error("cannot handle agent vers: ", err)
		return
	}

	lastAgentVers = digest
}

func handleAgVers(arr []AgVer) error {
	if arr == nil {
		return allUninstall()
	}

	cnt := len(arr)
	if cnt == 0 {
		return allUninstall()
	}

	remote := make(map[string]struct{}, cnt)

	for i := 0; i < cnt; i++ {
		remote[arr[i].Name] = struct{}{}

		ag, has := agentStats[arr[i].Name]
		if has && ag.Ver == arr[i].Ver {
			continue
		}

		if !has {
			err := agentInstall(arr[i].Name, arr[i].Ver)
			if err != nil {
				logger.Errorf("cannot install agent:%s-%s: %v", arr[i].Name, arr[i].Ver, err)
				return err
			}
			continue
		}

		// version not equal
		err := agentUpgrade(arr[i].Name, arr[i].Ver)
		if err != nil {
			logger.Errorf("cannot upgrade agent:%s, %s->%s, error: %v", arr[i].Name, ag.Ver, arr[i].Ver, err)
			return err
		}
	}

	for name := range agentStats {
		if _, has := remote[name]; !has {
			err := agentUninstall(name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func allUninstall() error {
	for name := range agentStats {
		err := agentUninstall(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func agentVersDigest(res HeartbeatResponse) string {
	if res.Ags == nil {
		return ""
	}

	cnt := len(res.Ags)
	if cnt == 0 {
		return ""
	}

	sort.Sort(AgVerSlice(res.Ags))

	arr := make([]string, 0, cnt*2)
	for i := 0; i < cnt; i++ {
		arr = append(arr, res.Ags[i].Name, res.Ags[i].Ver)
	}

	return strings.Join(arr, "")
}

func upgradeSgd(ver string) {
	tarFile := fmt.Sprintf("sgd-%s_%s_%s.tar.gz", ver, runtime.GOOS, runtime.GOARCH)
	md5File := fmt.Sprintf("sgd-%s_%s_%s.tar.gz.md5", ver, runtime.GOOS, runtime.GOARCH)
	tarPath := path.Join(config.TarDir, tarFile)
	md5Path := path.Join(config.TarDir, md5File)
	tarUrl := fmt.Sprintf("%s/tarball/%s", config.Srv, tarFile)
	md5Url := fmt.Sprintf("%s/tarball/%s", config.Srv, md5File)

	if err := ensureFilesReady(tarPath, md5Path, tarUrl, md5Url); err != nil {
		logger.Error(err)
		return
	}

	if err := sys.CmdRun("/bin/bash", "-c", fmt.Sprintf("tar zxf %s -C %s", tarPath, config.Cwd)); err != nil {
		logger.Error(err)
		return
	}

	// after finishing the arrangements, I can leave with ease
	os.Exit(0)
}

// do not trust local files
func ensureFilesReady(tarPath, md5Path, tarUrl, md5Url string) error {
	if err := file.Unlink(tarPath); err != nil {
		return fmt.Errorf("cannot rm %s: %v", tarPath, err)
	}

	if err := file.Unlink(md5Path); err != nil {
		return fmt.Errorf("cannot rm %s: %v", md5Path, err)
	}

	if err := file.Download(tarPath, tarUrl); err != nil {
		return fmt.Errorf("download %s to %s fail: %v", tarUrl, tarPath, err)
	}

	if err := file.Download(md5Path, md5Url); err != nil {
		return fmt.Errorf("download %s to %s fail: %v", md5Url, md5Path, err)
	}

	succ, err := md5c(tarPath, md5Path)
	if err != nil {
		return fmt.Errorf("cannot exec md5c: %v", err)
	}

	if !succ {
		return fmt.Errorf("md5sum check not equal, tar:%s, md5:%s", tarPath, md5Path)
	}

	return nil
}

func md5c(tarPath, md5Path string) (bool, error) {
	md5content, err := file.ReadString(md5Path)
	if err != nil {
		return false, err
	}

	md5compute, err := file.Md5File(tarPath)
	if err != nil {
		return false, err
	}

	if strings.Contains(md5content, md5compute) {
		return true, nil
	}

	return false, nil
}
