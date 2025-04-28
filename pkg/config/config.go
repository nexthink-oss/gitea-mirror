package config

import (
	"fmt"
	"iter"
	"path/filepath"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

// Source may be either GitHub or Gitea
type Source struct {
	Type      string `mapstructure:"type" yaml:"type,omitempty" default:"gitea"`
	Url       string `mapstructure:"url" yaml:"url,omitempty"`
	RemoteUrl string `mapstructure:"remote-url" yaml:"remote-url,omitempty"`
	Token     string `mapstructure:"token" yaml:"-"`
}

// Target is always a Gitea instance
type Target struct {
	Url   string `mapstructure:"url" yaml:"url" default:"http://localhost:3000"`
	Token string `mapstructure:"token" yaml:"-"`
}

type Forge interface {
	GetType() string
	GetUrl() string
	GetRemoteUrl() string
	GetToken() string
}

type Defaults struct {
	Owner        string        `mapstructure:"owner"`
	Interval     time.Duration `mapstructure:"interval"`
	PublicSource bool          `mapstructure:"public-source"`
	PublicTarget bool          `mapstructure:"public-target"`
}

type Repository struct {
	Owner        string         `mapstructure:"owner" yaml:"owner,omitempty"`
	Name         string         `mapstructure:"name"`
	Interval     *time.Duration `mapstructure:"interval" yaml:"interval,omitempty"`
	PublicSource *bool          `mapstructure:"public-source" yaml:"public-source,omitempty"`
	PublicTarget *bool          `mapstructure:"public-target" yaml:"public-target,omitempty"`
}

type RepositorySet map[string]struct{}

type Config struct {
	Source       Source       `mapstructure:"source"`
	Target       Target       `mapstructure:"target"`
	Defaults     Defaults     `mapstructure:"defaults"`
	Repositories []Repository `mapstructure:"repositories"`
}

// LoadConfig loads the configuration using viper from the given file
// and returns the configuration object.
func LoadConfig(names []string) (*Config, error) {
	var config Config

	for _, name := range names {

		if extension := filepath.Ext(name); extension != "" {
			viper.SetConfigType(extension[1:])
		} else {
			viper.SetConfigType("yaml")
		}
		viper.SetConfigFile(name)

		if err := viper.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	for i, r := range config.Repositories {
		if r.Owner == "" {
			config.Repositories[i].Owner = config.Defaults.Owner
		}
		if r.Interval == nil {
			config.Repositories[i].Interval = &config.Defaults.Interval
		}
		if r.PublicSource == nil {
			config.Repositories[i].PublicSource = &config.Defaults.PublicSource
		}
		if r.PublicTarget == nil {
			config.Repositories[i].PublicTarget = &config.Defaults.PublicTarget
		}
	}

	return &config, nil
}

func (s *Source) GetType() string {
	return s.Type
}

func (s *Source) GetUrl() string {
	return s.Url
}

func (s *Source) GetRemoteUrl() string {
	return s.RemoteUrl
}

func (s *Source) GetToken() string {
	return s.Token
}

func (t *Target) GetType() string {
	return "gitea"
}

func (t *Target) GetUrl() string {
	return t.Url
}

func (t *Target) GetRemoteUrl() string {
	return ""
}

func (t *Target) GetToken() string {
	return t.Token
}

func (r *Repository) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

func (r *Repository) Success(message ...string) string {
	if len(message) > 0 {
		return fmt.Sprintf("✅ %s/%s: %s", r.Owner, r.Name, strings.Join(message, " "))
	}
	return fmt.Sprintf("✅ %s/%s", r.Owner, r.Name)
}

func (r *Repository) Failure(err error) string {
	return fmt.Sprintf("❌ %s/%s: %s", r.Owner, r.Name, err.Error())
}

func (c *Config) RepositorySetFromArgs(args []string) RepositorySet {
	sets := make(RepositorySet)
	for _, arg := range args {
		repo, err := c.ParseRepositorySpec(arg)
		if err != nil {
			continue
		}
		sets[repo] = struct{}{}
	}
	return sets
}

func (c *Config) ParseRepositorySpec(spec string) (repo string, err error) {
	split := strings.SplitN(spec, "/", 2)
	switch len(split) {
	case 1:
		repo = fmt.Sprintf("%s/%s", c.Defaults.Owner, split[0])
	case 2:
		repo = fmt.Sprintf("%s/%s", split[0], split[1])
	default:
		err = fmt.Errorf("invalid repository spec")
	}
	return repo, err
}

func (c *Config) FilteredRepositories(args []string) iter.Seq[Repository] {
	filter := c.RepositorySetFromArgs(args)
	return func(yield func(Repository) bool) {
		for _, repo := range c.Repositories {
			if filter.EmptyOrContains(repo) && !yield(repo) {
				return
			}
		}
	}
}

func (r *RepositorySet) EmptyOrContains(repo Repository) bool {
	if r == nil || len(*r) == 0 {
		return true
	}

	_, ok := (*r)[repo.String()]
	return ok
}
