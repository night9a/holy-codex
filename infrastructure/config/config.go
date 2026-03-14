package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	//"github.com/gogo/protobuf/plugin/defaultcheck"
)

type Config struct {
	DBPath      string `json:"db_path"`
	NetworkPort int    `json:"network_port"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	AutoSaveMs  int    `json:"autosave_ms"`
	SyncEnabled bool   `json:"sync_enabled"`
}

const configFile = "holydiary.json"

func Load() (*Config,error) {
	dir ,err := configDir()
	if err != nil {
		return nil,err
	}
	path := filepath.Join(dir,configFile)

	cfg := defaults(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, cfg.Save()
		}
		return nil,err
	}

	if err := json.Unmarshal(data,cfg); err!= nil {
		return nil, err
	}
	return cfg,nil
}

func (c *Config) Save() error {
	dir ,err := configDir()
	if err != nil {
		return err
	}
	data ,err := json.MarshalIndent(c,""," ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir,configFile),data,0600)
}

func defaults(dir string) *Config {
	return &Config{
		DBPath:      filepath.Join(dir, "diary.db"),
		NetworkPort: 7654,
		UserID:      "default",
		UserName:    "Scribe",
		AutoSaveMs:  3000,
		SyncEnabled: true,
	}
}

func configDir() (string,error) {
	base,err := os.UserConfigDir()
	if err != nil {
		base,err = os.UserHomeDir()
		if err != nil {
			return "",err
		}
	}
	dir := filepath.Join(base,"HolyDiary")
	return dir, os.MkdirAll(dir, 0700)
}