package config

const (
	Ver = "1.1.5"
	Url = "UUID_URL"

	LogFileSize = 1024 * 1024 * 5
	LogFileNum  = 3

	Started = 1
	Stopped = 0
)

var (
	Dev  = false
	Srv  = ""
	UUID = ""
	Cwd  = ""

	LogDir = ""
	AgsDir = ""
	TarDir = ""
)
