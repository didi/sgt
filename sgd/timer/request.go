package timer

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"sgt/pkg/file"
	"sgt/sgd/config"
)

func makeHeartbeatRequest() (HeartbeatRequest, error) {
	req := HeartbeatRequest{
		UUID:   config.UUID,
		SgdVer: config.Ver,
	}

	err := updateAgentStats()
	if err != nil {
		return req, err
	}

	cnt := len(agentStats)
	ags := make([]AgSt, 0, cnt)
	for name, stat := range agentStats {
		ags = append(ags, AgSt{
			Name: name,
			Ver:  stat.Ver,
			Stat: stat.Stat,
		})
	}

	req.Ags = ags

	return req, nil
}

func updateAgentStats() error {
	agentDirs, err := file.DirsUnder(config.AgsDir)
	if err != nil {
		return fmt.Errorf("cannot list %s, error: %v", config.AgsDir, err)
	}

	sort.Strings(agentDirs)

	newAgentDirs := strings.Join(agentDirs, "")
	if lastAgentDirs != newAgentDirs {
		if err = collectAgents(agentDirs); err != nil {
			return err
		}
		lastAgentDirs = newAgentDirs
		return nil
	}

	changed := false
	for name := range agentStats {
		// if agent stopped, start it
		if agentStats[name].Stat == config.Stopped {
			agentStart(name)
			changed = true
		}
	}

	if changed {
		return collectAgents(agentDirs)
	}

	changed = false
	for name := range agentStats {
		cmdline, err := file.ReadString(fmt.Sprintf("/proc/%d/cmdline", agentStats[name].Pid))
		if err != nil {
			if config.Dev {
				log.Printf("INF: agent[%s] cannot get cmdline, error: %v", name, err)
			}
			changed = true
			break
		}

		if cmdline != agentStats[name].Cmdline {
			if config.Dev {
				log.Printf("INF: agent[%s] cmdline changed, old: %s, new: %s", agentStats[name].Cmdline, cmdline)
			}
			changed = true
			break
		}
	}

	if changed {
		return collectAgents(agentDirs)
	}

	return nil
}
