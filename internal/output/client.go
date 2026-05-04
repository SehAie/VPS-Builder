package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SehAie/VPS-Builder/internal/config"
)

func WriteHysteriaClient(cfg config.Config, password string) (string, error) {
	if err := os.MkdirAll(cfg.Paths.ClientDir, 0o700); err != nil {
		return "", err
	}
	path := cfg.ClientConfigPath()
	body := fmt.Sprintf(`server: %s:%d
auth: %q

tls:
  # self-signed mode needs insecure=true unless you pin/import the cert.
  insecure: true

socks5:
  listen: 127.0.0.1:1080

http:
  listen: 127.0.0.1:8080
`, cfg.Server.IP, cfg.Hysteria.ListenPort, password)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		return "", err
	}
	abs, _ := filepath.Abs(path)
	return abs, nil
}

func WriteSummary(cfg config.Config, password string) error {
	if err := os.MkdirAll(cfg.Paths.BaseDir, 0o700); err != nil {
		return err
	}
	path := filepath.Join(cfg.Paths.BaseDir, "last-deploy-summary.txt")
	body := fmt.Sprintf(`VPS-Builder deployment summary

Server IP: %s
SSH user: %s
SSH port: %d
SSH private key: %s
Hysteria 2: %v
Hysteria UDP port: %d
Hysteria password: %s
Client config: %s

Important:
- Back up the private key and client config.
- Make sure your VPS provider security group allows TCP %d and UDP %d.
`, cfg.Server.IP, cfg.Server.User, cfg.SSH.NewPort, cfg.KeyPath(), cfg.Hysteria.Enabled, cfg.Hysteria.ListenPort, redact(password), cfg.ClientConfigPath(), cfg.SSH.NewPort, cfg.Hysteria.ListenPort)
	return os.WriteFile(path, []byte(body), 0o600)
}

func redact(s string) string {
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}
