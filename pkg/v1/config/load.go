package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	SmtpC Smtp
	LdapC Ldap
)

func LoadConfig() error {

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	var cfg struct {
		SMTP Smtp `yaml:"smtp"`
		LDAP Ldap `yaml:"ldap"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	SmtpC = cfg.SMTP
	LdapC = cfg.LDAP

	return nil
}
