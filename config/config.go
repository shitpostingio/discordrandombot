package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

// Config is the bot configuration representation, read
// from a configuration file.
type Config struct {
	DiscordTokenBot string
}

// ReadConfig loads the values from the config file
func ReadConfig(path string) (Config, error) {
	var conf Config

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return Config{}, err
	}

	if conf.DiscordTokenBot == "" {
		return newErr("missing Bot token")
	}

	return conf, nil
}

func newErr(message string) (Config, error) {
	return Config{}, errors.New(message)
}
