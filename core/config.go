package core

import (
	json "encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	HttpAddr string `json:"httpAddr"`
	RaftAddr string `json:"raftAddr"`
}

func InitConfig() (*Config, error) {
	f, err := ioutil.ReadFile("core/config.json")
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal([]byte(f), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func InitLogger() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	log.SetLevel(log.DebugLevel)
}

func GetLogger(name ...string) *log.Entry {
	if len(name) > 0 {
		return log.WithField("module", name[0])
	}
	return log.WithField("module", nil)
}
