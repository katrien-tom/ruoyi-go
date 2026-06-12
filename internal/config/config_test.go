package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	validConfig := `mysql:
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: test
  charset: utf8mb4
  max-idle-conns: 10
  max-open-conns: 100
  conn-max-lifetime: 1h

redis:
  host: 127.0.0.1
  port: 6379
  username: ""
  password: ""
  db: 0
  pool-size: 10
  min-idle-conns: 5
  dial-timeout: 5s
  read-timeout: 3s
  write-timeout: 3s
`

	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.MySQL.Host != "127.0.0.1" {
		t.Errorf("expected mysql host 127.0.0.1, got %s", cfg.MySQL.Host)
	}
	if cfg.MySQL.Port != 3306 {
		t.Errorf("expected mysql port 3306, got %d", cfg.MySQL.Port)
	}
	if cfg.Redis.Host != "127.0.0.1" {
		t.Errorf("expected redis host 127.0.0.1, got %s", cfg.Redis.Host)
	}
}

func TestLoadConfigValidationFails(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	invalidConfig := `mysql:
  host: ""
  port: 99999
  username: root
  database: test
  charset: utf8mb4

redis:
  host: 127.0.0.1
  port: 6379
  db: 0
  pool-size: 10
  min-idle-conns: 5
  dial-timeout: 5s
  read-timeout: 3s
  write-timeout: 3s
`

	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
