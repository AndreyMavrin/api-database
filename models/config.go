package models

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type SQLDataBase struct {
	Server          string   `toml:"Server"`
	Database        string   `toml:"Database"`
	ApplicationName string   `toml:"ApplicationName"`
	MaxIdleConns    int      `toml:"MaxIdleConns"`
	MaxOpenConns    int      `toml:"MaxOpenConns"`
	ConnMaxLifetime duration `toml:"ConnMaxLifetime"`
	UserID          string
	Password        string
}

type duration time.Duration

var (
	configPath = "config/config.toml"
	hashPaths  = []string{configPath}
)

type Config struct {
	SQLDataBase SQLDataBase `toml:"SQLDataBase"`
}

func (d *duration) UnmarshalText(text []byte) error {
	temp, err := time.ParseDuration(string(text))
	*d = duration(temp)
	return err
}

func LoadConfig(c *Config) {
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return
	}
	c.SQLDataBase.UserID = getCredential("/etc/scrt/forums/sqlUser")
	c.SQLDataBase.Password = getCredential("/etc/scrt/forums/sqlPassword")
}

func getCredential(path string) string {
	hashPaths = append(hashPaths, path)
	c, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	return strings.TrimSpace(string(c))
}
