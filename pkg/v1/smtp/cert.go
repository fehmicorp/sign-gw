package smtp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/logger"
	"go.uber.org/zap"
)

func GenerateCertificate() error {

	parts := strings.Split(
		config.SmtpC.RelayUsername,
		"@",
	)

	if len(parts) != 2 {
		return fmt.Errorf("invalid relay username")
	}

	host := "mail." + strings.ToLower(parts[1])

	certDir := "./certs"

	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	if err := os.MkdirAll(certDir, 0755); err != nil {
		return err
	}

	// Certificate already exists
	if _, err := os.Stat(certFile); err == nil {

		if _, err := os.Stat(keyFile); err == nil {

			logger.Info(
				"TLS certificate already exists",
				zap.String("host", host),
				zap.String("cert", certFile),
			)

			return nil
		}
	}

	mkcert := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"Microsoft",
		"WinGet",
		"Packages",
		"FiloSottile.mkcert_Microsoft.Winget.Source_8wekyb3d8bbwe",
		"mkcert.exe",
	)

	logger.Info(
		"Generating TLS certificate",
		zap.String("host", host),
	)

	cmd := exec.Command(
		mkcert,
		"-cert-file", certFile,
		"-key-file", keyFile,
		host,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	logger.Info(
		"TLS certificate generated",
		zap.String("host", host),
		zap.String("cert", certFile),
		zap.String("key", keyFile),
	)

	return nil
}
