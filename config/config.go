package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DefaultFileName string
	TimeoutSeconds  string
}

func (c *Config) loadConfigFile(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return err
	}
	c.DefaultFileName = os.Getenv("FILENAME")
	c.TimeoutSeconds = os.Getenv("TIMEOUTSECONDS")
	return nil
}

func LoadConfig(file string) (Config, error) {
	var c Config
	err := c.loadConfigFile(file)
	if err != nil {
		return c, err
	}
	return c, nil
}
