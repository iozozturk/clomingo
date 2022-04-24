package config

import (
	"embed"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"os"
)

//go:embed config_dev.json config_local.json config_prod.json
var configs embed.FS

type Config interface {
}

type ApplicationConf struct {
	Environment     string `envconfig:"ENV"`
	Host            string
	GoogleClientIds []string
	GoogleProjectId string
}

func NewApplication() *ApplicationConf {
	env := os.Getenv("ENV")
	var confPath string
	if env == "" || env == "local" {
		confPath = "config_local.json"
	} else {
		switch env {
		case "dev":
			confPath = "config_dev.json"
		case "prod":
			confPath = "config_prod.json"
		}
	}
	file, fileErr := configs.Open(confPath)
	if fileErr != nil {
		panic(fileErr)
	}
	decoder := json.NewDecoder(file)
	var conf ApplicationConf
	err := decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	err = envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
