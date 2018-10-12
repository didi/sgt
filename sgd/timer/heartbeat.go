package timer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"sgt/pkg/logger"
	"sgt/sgd/config"
)

func Heartbeat() {
	dur := time.Duration(3) * time.Second
	url := fmt.Sprintf("%s/heartbeat", config.Srv)
	for {
		heartbeat(url)
		time.Sleep(dur)
	}
}

func heartbeat(url string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			return
		}
	}()

	req, err := makeHeartbeatRequest()
	if err != nil {
		logger.Error("cannot make heartbeat request: ", err)
		return
	}

	bs, err := json.Marshal(req)
	if err != nil {
		logger.Error("cannot marshal heartbeat request: ", err)
		return
	}

	if config.Dev {
		log.Println("INF: heartbeat request ->", string(bs))
	}

	res, err := httpcli.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		logger.Error("cannot dial heartbeat server: ", err)
		return
	}

	if res.StatusCode != 200 {
		logger.Error("heartbeat server return statuscode: ", res.StatusCode)
		return
	}

	if res.Body == nil {
		logger.Error("heartbeat server response body is nil")
		return
	}

	defer res.Body.Close()

	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("cannot read heartbeat server response: ", err)
		return
	}

	if config.Dev {
		log.Println("INF: heartbeat response ->", string(bs))
	}

	var reply HeartbeatResponse
	err = json.Unmarshal(bs, &reply)
	if err != nil {
		logger.Error("cannot unmarshal heartbeat server response: ", err)
		return
	}

	handleHeartbeatResponse(reply)

	checkMem()
}
