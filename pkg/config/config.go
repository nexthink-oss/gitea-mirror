package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

// Source may be either GitHub or Gitea
type Source struct {
	Type  string `mapstructure:"type" yaml:"type,omitempty" default:"gitea"`
	Url   string `mapstructure:"url" yaml:"url,omitempty"`
	Token string `mapstructure:"token" yaml:"-"`
}

// Target is always a Gitea instance
type Target struct {
	Url   string `mapstructure:"url" yaml:"url" default:"http://localhost:3000"`
	Token string `mapstructure:"token" yaml:"-"`
}

type Defaults struct {
	Owner    string        `mapstructure:"owner"`
	Interval time.Duration `mapstructure:"interval"`
	Public   bool          `mapstructure:"public"`
}

type Repository struct {
	Owner    *string        `mapstructure:"owner" yaml:"owner,omitempty"`
	Name     string         `mapstructure:"name"`
	Interval *time.Duration `mapstructure:"interval" yaml:"interval,omitempty"`
	Public   *bool          `mapstructure:"public" yaml:"public,omitempty"`
}

type Config struct {
	Source       Source       `mapstructure:"source"`
	Target       Target       `mapstructure:"target"`
	Defaults     Defaults     `mapstructure:"defaults"`
	Repositories []Repository `mapstructure:"repositories"`
}

// LoadConfig loads the configuration using viper from the given file
// and returns the configuration object.
func LoadConfig(name string, paths ...string) (*Config, error) {
	var config Config

	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigName(name)

	viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	for i, r := range config.Repositories {
		if r.Owner == nil {
			config.Repositories[i].Owner = &config.Defaults.Owner
		}
		if r.Interval == nil {
			config.Repositories[i].Interval = &config.Defaults.Interval
		}
		if r.Public == nil {
			config.Repositories[i].Public = &config.Defaults.Public
		}
	}

	return &config, nil
}

func (r Repository) Success(message ...string) string {
	if len(message) > 0 {
		return fmt.Sprintf("✅ %s/%s: %s", *r.Owner, r.Name, strings.Join(message, " "))
	}
	return fmt.Sprintf("✅ %s/%s", *r.Owner, r.Name)
}

func (r Repository) Failure(err error) string {
	return fmt.Sprintf("❌ %s/%s: %s", *r.Owner, r.Name, err.Error())
}
