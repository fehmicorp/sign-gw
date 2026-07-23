package logger

type Config struct {
	Level      string `json:"level" yaml:"level"`
	Console    bool   `json:"console" yaml:"console"`
	File       string `json:"file" yaml:"file"`
	MaxSize    int    `json:"maxSize" yaml:"maxSize"`
	MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
	MaxAge     int    `json:"maxAge" yaml:"maxAge"`
	Compress   bool   `json:"compress" yaml:"compress"`
	Format     string `json:"format" yaml:"format"`
	Caller     bool   `json:"caller" yaml:"caller"`
	Stacktrace bool   `json:"stacktrace" yaml:"stacktrace"`
}
