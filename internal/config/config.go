package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName string         `yaml:"project_name"`
	Server      ServerConfig   `yaml:"server"`
	SSH         SSHConfig      `yaml:"ssh"`
	System      SystemConfig   `yaml:"system"`
	Hysteria    HysteriaConfig `yaml:"hysteria"`
	Paths       PathsConfig    `yaml:"paths"`
}

type ServerConfig struct {
	IP       string `yaml:"ip"`
	User     string `yaml:"user"`
	OldPort  int    `yaml:"old_port"`
	Password string `yaml:"password,omitempty"` // Not recommended. Prefer prompt or env VAKIT_ROOT_PASSWORD.
}

type SSHConfig struct {
	NewPort              int    `yaml:"new_port"`
	KeyName              string `yaml:"key_name"`
	DisablePasswordLogin bool   `yaml:"disable_password_login"`
	PermitRootLogin      string `yaml:"permit_root_login"`
	MaxAuthTries         int    `yaml:"max_auth_tries"`
}

type SystemConfig struct {
	UpdatePackages bool `yaml:"update_packages"`
	EnableBBR      bool `yaml:"enable_bbr"`
	OptimizeSysctl bool `yaml:"optimize_sysctl"`
	DisableIPv6    bool `yaml:"disable_ipv6"`
}

type HysteriaConfig struct {
	Enabled              bool      `yaml:"enabled"`
	ListenPort           int       `yaml:"listen_port"`
	Password             string    `yaml:"password"` // "auto" means generate on deploy.
	TLS                  TLSConfig `yaml:"tls"`
	GenerateClientConfig bool      `yaml:"generate_client_config"`
}

type TLSConfig struct {
	Mode       string `yaml:"mode"` // self-signed or acme. MVP deploys self-signed only.
	CommonName string `yaml:"common_name"`
	Domain     string `yaml:"domain,omitempty"`
	Email      string `yaml:"email,omitempty"`
}

type PathsConfig struct {
	BaseDir   string `yaml:"base_dir"`
	KeysDir   string `yaml:"keys_dir"`
	ClientDir string `yaml:"client_dir"`
}

func Default() Config {
	base := defaultBaseDir()
	return Config{
		ProjectName: "vps-auto-kit-windows",
		Server: ServerConfig{
			User:    "root",
			OldPort: 22,
		},
		SSH: SSHConfig{
			NewPort:              2222,
			KeyName:              "vps-auto-key",
			DisablePasswordLogin: true,
			PermitRootLogin:      "prohibit-password",
			MaxAuthTries:         3,
		},
		System: SystemConfig{
			UpdatePackages: true,
			EnableBBR:      true,
			OptimizeSysctl: true,
			DisableIPv6:    false,
		},
		Hysteria: HysteriaConfig{
			Enabled:              true,
			ListenPort:           8443,
			Password:             "auto",
			GenerateClientConfig: true,
			TLS: TLSConfig{
				Mode:       "self-signed",
				CommonName: "bing.com",
			},
		},
		Paths: PathsConfig{
			BaseDir:   base,
			KeysDir:   filepath.Join(base, "keys"),
			ClientDir: filepath.Join(base, "clients"),
		},
	}
}

func defaultBaseDir() string {
	if runtime.GOOS == "windows" {
		if v := os.Getenv("APPDATA"); v != "" {
			return filepath.Join(v, "vps-auto-kit")
		}
	}
	if v := os.Getenv("HOME"); v != "" {
		return filepath.Join(v, ".vps-auto-kit")
	}
	return ".vps-auto-kit"
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	cfg := Default()
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Paths.BaseDir == "" {
		cfg.Paths.BaseDir = defaultBaseDir()
	}
	if cfg.Paths.KeysDir == "" {
		cfg.Paths.KeysDir = filepath.Join(cfg.Paths.BaseDir, "keys")
	}
	if cfg.Paths.ClientDir == "" {
		cfg.Paths.ClientDir = filepath.Join(cfg.Paths.BaseDir, "clients")
	}
	return cfg, cfg.Validate()
}

func Save(path string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, b, 0o600)
}

func (c Config) Validate() error {
	var problems []string
	if strings.TrimSpace(c.Server.IP) == "" || net.ParseIP(c.Server.IP) == nil {
		problems = append(problems, "server.ip must be a valid IP address")
	}
	if strings.TrimSpace(c.Server.User) == "" {
		problems = append(problems, "server.user is required")
	}
	if c.Server.OldPort < 1 || c.Server.OldPort > 65535 {
		problems = append(problems, "server.old_port must be 1-65535")
	}
	if c.SSH.NewPort < 1 || c.SSH.NewPort > 65535 {
		problems = append(problems, "ssh.new_port must be 1-65535")
	}
	if c.Hysteria.ListenPort < 1 || c.Hysteria.ListenPort > 65535 {
		problems = append(problems, "hysteria.listen_port must be 1-65535")
	}
	if c.SSH.KeyName == "" {
		problems = append(problems, "ssh.key_name is required")
	}
	if c.SSH.PermitRootLogin == "" {
		problems = append(problems, "ssh.permit_root_login is required")
	}
	if c.Hysteria.TLS.Mode != "self-signed" && c.Hysteria.TLS.Mode != "acme" {
		problems = append(problems, "hysteria.tls.mode must be self-signed or acme")
	}
	if c.Hysteria.TLS.Mode == "acme" {
		problems = append(problems, "ACME mode is scaffolded but not deployed in this MVP; use self-signed")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

func (c Config) KeyPath() string {
	cleanIP := strings.ReplaceAll(c.Server.IP, ":", "_")
	return filepath.Join(c.Paths.KeysDir, fmt.Sprintf("%s_%s", cleanIP, c.SSH.KeyName))
}

func (c Config) ClientConfigPath() string {
	cleanIP := strings.ReplaceAll(c.Server.IP, ":", "_")
	return filepath.Join(c.Paths.ClientDir, fmt.Sprintf("%s_hysteria.yaml", cleanIP))
}
