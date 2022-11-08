package config

import (
	"github.com/hg/pgstaging/consts"
	"github.com/hg/pgstaging/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"os/user"
	"path"
)

type Config struct {
	// Which address to listen on.
	Listen string `json:"listen"`

	// Administrator password with access to all databases.
	Passwd string `json:"passwd"`

	// Web server drops privileges from root to this user.
	User string `json:"user"`
}

func (c *Config) validate() error {
	if _, err := netip.ParseAddrPort(c.Listen); err != nil {
		return err
	}
	if _, err := user.Lookup(c.User); err != nil {
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
			conf := defaultConfig()
			trySaveConfig(cp, conf)
			return conf, nil
		}
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	conf := Config{}

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	if conf.Passwd == "" {
		conf.Passwd = util.RandomString(18)
		log.Printf("replacing empty admin password with a one-time random one: '%s'", conf.Passwd)
	}

	return &conf, conf.validate()
}

func trySaveConfig(p string, conf *Config) {
	err := os.MkdirAll(path.Base(p), 0o755)
	if err != nil {
		log.Printf("could not create config dir: %v", err)
		return
	}
	b, err := json.Marshal(conf)
	if err != nil {
		log.Printf("could not marshal config: %v", err)
		return
	}
	err = os.WriteFile(p, b, 0o600)
	if err != nil {
		log.Printf("could not write config: %v", err)
	}
}
