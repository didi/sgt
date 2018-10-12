package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"sgt/pkg/file"
)

func Init() {
	Cwd = file.Dir(file.SelfDir())
	if Dev {
		log.Printf("INF: workdir -> %s", Cwd)
		log.Printf("INF: version -> %s", Ver)
	}

	LogDir = path.Join(Cwd, "log")
	AgsDir = path.Join(Cwd, "ags")
	TarDir = path.Join(Cwd, "tar")

	uuid, err := GetUUID()
	if err == nil {
		SetUUID(uuid)
		return
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		uuid, err = GetUUID()
		if err == nil {
			SetUUID(uuid)
			return
		}
	}

	log.Fatalf("FAT: cannot get uuid")
}

func SetUUID(uuid string) {
	if Dev {
		log.Printf("INF: uuid -> %s", uuid)
	}
	UUID = uuid
}

func GetUUID() (string, error) {
	client := http.Client{
		Timeout: time.Second,
	}

	res, err := client.Get(Url)
	if err != nil {
		return "", fmt.Errorf("cannot dial uuid url: %v", err)
	}

	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read uuid body: %v", err)
	}

	return string(bs), nil
}
