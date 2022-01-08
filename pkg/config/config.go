package config

import (
	"encoding/json"
	"os"
)

type Config interface {
}

type ApplicationConf struct {
	Host            string
	GoogleClientIds []string
	GoogleProjectId string
}

func NewApplication() *ApplicationConf {
	env := os.Getenv("ENV")
	var confPath string
	if env == "" {
		confPath = "pkg/config/config_dev.json"
	} else {
		confPath = "pkg/config/config_prod.json"
	}
	file, fileErr := os.Open(confPath)
	if fileErr != nil {
		panic(fileErr)
	}
	decoder := json.NewDecoder(file)
	var conf ApplicationConf
	err := decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
