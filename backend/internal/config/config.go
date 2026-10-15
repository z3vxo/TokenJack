package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	TLS     TLSConfig     `toml:"TLS"`
	Auth    AuthConfig    `toml:"auth"`
	Logging LoggingConfig `toml:"logging"`
}

type TLSConfig struct {
	Enabled  bool   `toml:"enabled"`
	Hostname string `toml:"hostname"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}
type LoggingConfig struct {
	Enabled bool   `toml:"enabled"`
	Level   string `toml:"level"`
}

type AuthConfig struct {
	AccessSecret  string        `toml:"access_secret"`
	RefreshSecret string        `toml:"refresh_secret"`
	AccessTTL     time.Duration `toml:"access_ttl"`
	RefreshTTL    time.Duration `toml:"refresh_ttl"`
}

type ServerConfig struct {
	Bind string `toml:"bind"`
	Port int    `toml:"port"`
}

func Generate(path string) error {
	var cfg Config
	cfg.Server.Bind = "0.0.0.0"
	cfg.Server.Port = 8080

	cfg.TLS.Enabled = false
	cfg.TLS.Hostname = "127.0.0.1"
	cfg.TLS.CertFile = ""
	cfg.TLS.KeyFile = ""

	accessSecret, err := genSecret()
	if err != nil {
		return err
	}
	refreshSecret, err := genSecret()
	if err != nil {
		return err
	}
	cfg.Auth.AccessSecret = accessSecret
	cfg.Auth.RefreshSecret = refreshSecret
	cfg.Auth.AccessTTL = 15 * time.Minute
	cfg.Auth.RefreshTTL = 7 * 24 * time.Hour

	cfg.Logging.Enabled = true
	cfg.Logging.Level = "info"

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if _, err := f.WriteString("# leave tls key and cert path empty for server to generate self signed keys.\n"); err != nil {
		return err
	}

	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}

func genSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Get() (*Config, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfgPath := filepath.Join(h, ".local", "share", "token-jack", "config", "jacker.conf")

	var cfg Config
	if _, err := toml.DecodeFile(cfgPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
