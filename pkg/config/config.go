package config

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

// Source may be either GitHub or Gitea
type Source struct {
	Type  string `mapstructure:"type" yaml:"type,omitempty" default:"gitea"`
	Url   string `mapstructure:"url" yaml:"url,omitempty"`
	Token string `mapstructure:"token"`
}

// Target is always a Gitea instance
type Target struct {
	Url   string `mapstructure:"url" yaml:"url" default:"http://localhost:3000"`
	Token string `mapstructure:"token" yaml:"token"`
}

type Repository struct {
	Owner   string `mapstructure:"owner" default:"nexthink"`
	Name    string `mapstructure:"name"`
	Private bool   `mapstructure:"private" yaml:"private,omitempty" default:"true"`
}

type Config struct {
	Source       Source       `json:"source"`
	Target       Target       `json:"target"`
	Repositories []Repository `json:"repositories"`
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

	return &config, nil
}

func (r Repository) Success(message ...string) string {
	if len(message) > 0 {
		return fmt.Sprintf("✅ %s/%s: %s", r.Owner, r.Name, strings.Join(message, " "))
	}
	return fmt.Sprintf("✅ %s/%s", r.Owner, r.Name)
}

func (r Repository) Failure(err error) string {
	return fmt.Sprintf("❌ %s/%s: %s", r.Owner, r.Name, err.Error())
}
