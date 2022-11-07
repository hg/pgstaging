package config

import (
	"github.com/hg/pgstaging/consts"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"path"
)

type Config struct {
	Listen string `json:"listen"`
}

func (c *Config) validate() error {
	_, err := netip.ParseAddrPort(c.Listen)
	if err != nil {
		return err
	}
	return nil
}

func defaultConfig() *Config {
	return &Config{
		Listen: ":80",
	}
}

func getConfigPath() string {
	dir := os.Getenv("CONFIGURATION_DIRECTORY")
	if dir == "" {
		dir = "/etc/" + consts.AppName
	}
	return path.Join(dir, "config.json")
}

func Load() (*Config, error) {
	cp := getConfigPath()
	log.Printf("using config %s", cp)

	bytes, err := os.ReadFile(cp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Print("config does not exist, using default")
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	conf := Config{}

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	return &conf, conf.validate()
}
