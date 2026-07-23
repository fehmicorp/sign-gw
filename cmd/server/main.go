package main

import (
	"log"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
	"github.com/fehmicorp/sign-gw/pkg/v1/smtp"
	"go.uber.org/zap"
)

func main() {
	config.Init()
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := logger.Init(config.Conf.Logging); err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
	defer logger.Sync()
	logger.Info(
		"Starting FEHMI Signature Gateway",
		zap.String("version", config.Conf.Application.Version),
		zap.String("company", config.Conf.Application.Company),
	)
	if err := config.Load(); err != nil {
		logger.Fatal("template load failed", zap.Error(err))
	}
	logger.Info("Templates Loaded", zap.Int("count", len(config.Templates)))
	config.StartTemplateWatcher()

	// conn, _ := ldap.Connect()
	// defer conn.Close()

	// ldap.Conn = conn

	logger.Info("Checking TLS Certificate")
	if err := smtp.GenerateCertificate(); err != nil {
		logger.Fatal(
			"certificate generation failed",
			zap.Error(err),
		)
	}
	logger.Info("TLS Certificate Ready")

	if err := smtp.Start(); err != nil {
		logger.Fatal("smtp server", zap.Error(err))
	}
	logger.Info(
		"Starting SMTP Server",
		zap.String("host", config.SmtpC.ListenHost),
		zap.Int("port", config.SmtpC.ListenPort),
		zap.String("hostname", config.SmtpC.Hostname),
		zap.Bool("tls", config.SmtpC.UseTLS),
	)
	// test()
	// testnew()

}
