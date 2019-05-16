package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

// Config configures the application.
type Config struct {
	HTTPPort    int             `yaml:"http_port" mapstructure:"http_port"`
	LogLevel    string          `yaml:"log_level" mapstructure:"log_level"`
	CORSEnabled bool            `yaml:"cors_enabled" mapstructure:"cors_enabled"`
	Database    *DatabaseConfig `yaml:"database" mapstructure:"database"`
	Github      *GithubConfig   `yaml:"github" mapstructure:"github"`
}

// DatabaseConfig configures the MySQL database connection.
type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Database string `yaml:"database" mapstructure:"database"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

// GithubConfig configures the Github authentication.
type GithubConfig struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
}

// LoadConfig loads a config from a YAML file.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)

	if err != nil {
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, err
	}

	var config Config
	var md mapstructure.Metadata
	msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Metadata:   &md,
		Result:     &config,
	})

	if err != nil {
		return nil, err
	}
	if err := msdec.Decode(raw); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate validates the config.
func (c *Config) Validate() error {
	if c.HTTPPort <= 0 {
		return fmt.Errorf("http_port is invalid")
	}
	if c.LogLevel == "" {
		return fmt.Errorf("log_level cannot be empty")
	}

	if c.Database == nil {
		return fmt.Errorf("database cannot be empty")
	}
	if err := c.Database.validate(); err != nil {
		return fmt.Errorf("database: %s", err)
	}

	if c.Github == nil {
		return fmt.Errorf("github cannot be empty")
	}
	if err := c.Github.validate(); err != nil {
		return fmt.Errorf("github: %s", err)
	}

	return nil
}

func (c *DatabaseConfig) validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Database == "" {
		return fmt.Errorf("database cannot be empty")
	}
	if c.User == "" {
		return fmt.Errorf("user cannot be empty")
	}
	if c.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

func (c *GithubConfig) validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("client_id cannot be empty")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client_secret cannot be empty")
	}

	return nil
}
