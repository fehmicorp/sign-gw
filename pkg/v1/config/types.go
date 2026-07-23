package config

import (
	"flag"

	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
)

type Config struct {
	Application ApplicationConfig `json:"application"`
	Logging     logger.Config     `json:"logging"`
	Template    string            `json:"template"`
}

type ApplicationConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Company string `json:"company"`
	Gitrepo string `json:"gitrepo"`
}

var ConfigFile string

func Init() {
	flag.StringVar(
		&ConfigFile,
		"c",
		"./config.yaml",
		"Configuration file",
	)
	flag.Parse()
}
