package timer

type AgentStat struct {
	Name    string
	Ver     string
	Pid     int
	Cmdline string
	Stat    int
}

// Stat: 0: stopped, 1: started
type AgSt struct {
	Name string `json:"n"`
	Ver  string `json:"v"`
	Stat int    `json:"s"`
}

type HeartbeatRequest struct {
	UUID   string `json:"id"`
	SgdVer string `json:"sgv"`
	Ags    []AgSt `json:"ags"`
}

type AgVer struct {
	Name string `json:"n"`
	Ver  string `json:"v"`
}

type HeartbeatResponse struct {
	Error  string  `json:"err"`
	SgdVer string  `json:"sgv"`
	Ags    []AgVer `json:"ags"`
}

type AgVerSlice []AgVer

func (s AgVerSlice) Len() int {
	return len(s)
}
func (s AgVerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s AgVerSlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

var (
	agentStats    = make(map[string]AgentStat)
	lastAgentDirs = "NULL"
	lastAgentVers = "NULL"
)
