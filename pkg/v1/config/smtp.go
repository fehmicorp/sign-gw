package config

import (
	"time"
)

type InlineImage struct {
	ContentID   string
	ContentType string
	Data        []byte
}

type Smtp struct {
	ListenHost string `yaml:"listenHost"`
	ListenPort int    `yaml:"listenPort"`
	Hostname   string `yaml:"hostname"`

	UseTLS        bool `yaml:"useTLS"`
	AllowInsecure bool `yaml:"allowInsecure"`

	MaxRecipients  int   `yaml:"maxRecipients"`
	MaxMessageSize int64 `yaml:"maxMessageSize"`

	SaveRawEML bool `yaml:"saveRawEML"`
	LogSMTP    bool `yaml:"logSMTP"`

	RelayHost     string `yaml:"relayHost"`
	RelayPort     int    `yaml:"relayPort"`
	RelayUsername string `yaml:"relayUsername"`
	RelayPassword string `yaml:"relayPassword"`
	RelayTLS      bool   `yaml:"relayTLS"`
}

type Attachment struct {
	FileName    string
	ContentType string
	ContentID   string
	Inline      bool
	Data        []byte
}

type Address struct {
	Name    string
	Address string
}

type Email struct {
	EnvelopeFrom string
	EnvelopeTo   []string
	From         Address
	To           []Address
	Cc           []Address
	Bcc          []Address
	ReplyTo      []Address
	Subject      string
	MessageID    string
	Date         time.Time
	Headers      map[string][]string
	Text         string
	HTML         string
	Raw          []byte
	InlineImages []InlineImage
	Attachments  []Attachment
}
