package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	SmtpC Smtp
	LdapC Ldap
	SaveC Save
)

type Save struct {
	Orignal bool `yaml:"orignal"`
	Edited  bool `yaml:"edited"`
}

func LoadConfig() error {

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	var cfg struct {
		SMTP Smtp `yaml:"smtp"`
		LDAP Ldap `yaml:"ldap"`
		SAVE Save `yaml:"save"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	SmtpC = cfg.SMTP
	LdapC = cfg.LDAP
	SaveC = cfg.SAVE

	return nil
}
