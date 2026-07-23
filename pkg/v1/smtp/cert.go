package smtp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	certDir := "./data/certs"

	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	if err := os.MkdirAll(certDir, 0755); err != nil {
		return err
	}

	// ------------------------------------------------------------
	// Already Exists
	// ------------------------------------------------------------

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

	logger.Info(
		"Generating TLS certificate",
		zap.String("host", host),
	)

	// ------------------------------------------------------------
	// Windows (mkcert)
	// ------------------------------------------------------------

	if runtime.GOOS == "windows" {

		mkcert := filepath.Join(
			os.Getenv("LOCALAPPDATA"),
			"Microsoft",
			"WinGet",
			"Packages",
			"FiloSottile.mkcert_Microsoft.Winget.Source_8wekyb3d8bbwe",
			"mkcert.exe",
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

	} else {

		// --------------------------------------------------------
		// Linux / Docker (OpenSSL)
		// --------------------------------------------------------

		cmd := exec.Command(
			"openssl",
			"req",
			"-x509",
			"-nodes",
			"-newkey", "rsa:4096",
			"-sha256",
			"-days", "3650",
			"-keyout", keyFile,
			"-out", certFile,
			"-subj", fmt.Sprintf("/CN=%s", host),
			"-addext", fmt.Sprintf("subjectAltName=DNS:%s", host),
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	logger.Info(
		"TLS certificate generated",
		zap.String("host", host),
		zap.String("cert", certFile),
		zap.String("key", keyFile),
	)

	return nil
}
