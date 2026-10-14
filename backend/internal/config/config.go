package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	TLS     TLSConfig     `toml:"TLS"`
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

type ServerConfig struct {
	Bind string `toml:"bind"`
	Port int    `tom:"port"`
}

func Generate(path string) error {
	var cfg Config
	cfg.Server.Bind = "0.0.0.0"
	cfg.Server.Port = 8080

	cfg.TLS.Enabled = false
	cfg.TLS.Hostname = "127.0.0.1"
	cfg.TLS.CertFile = ""
	cfg.TLS.KeyFile = ""

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
