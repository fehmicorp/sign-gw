package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
)

var Conf = Config{
	Application: ApplicationConfig{
		Name:    "FEHMI Signature Gateway",
		Version: "1.0.0",
		Company: "FEHMI Corporation",
		Gitrepo: "github.com/fehmicorp/sign-gw/v1",
	},
	Logging: logger.Config{
		Level:      "info",
		Console:    true,
		File:       "./logs/go_main.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Format:     "console",
		Caller:     true,
		Stacktrace: true,
	},
}

type Template struct {
	ID      string
	Name    string
	Default bool
	HTML    string
}

var Templates []Template

func StartTemplateWatcher() {

	go func() {

		for {

			if err := Load(); err != nil {

				log.Printf(
					"[Template] reload failed: %v",
					err,
				)

			} else {

				log.Printf(
					"[Template] loaded %d template(s)",
					len(Templates),
				)
			}

			time.Sleep(30 * time.Minute)
		}
	}()
}
func Load() error {
	dir := "./data/templates"
	Templates = Templates[:0]

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))

		if ext != ".html" && ext != ".htm" {
			continue
		}

		path := filepath.Join(dir, file.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		id := strings.TrimSuffix(file.Name(), ext)
		t := Template{
			ID:      id,
			Name:    strings.Title(strings.ReplaceAll(id, "_", " ")),
			Default: len(Templates) == 0,
			HTML:    string(data),
		}

		Templates = append(Templates, t)
	}
	return nil
}

func Get(id string) *Template {
	id = strings.ToLower(strings.TrimSpace(id))
	for i := range Templates {
		if strings.EqualFold(Templates[i].ID, id) {
			return &Templates[i]
		}
	}
	for i := range Templates {
		if Templates[i].Default {
			return &Templates[i]
		}
	}
	return nil
}
