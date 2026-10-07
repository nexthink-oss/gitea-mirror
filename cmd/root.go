package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/nexthink-oss/gitea-mirror/pkg/config"
	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/github"
	"github.com/nexthink-oss/gitea-mirror/pkg/server"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

var config *cfg.Config

func New() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "gitea-mirror",
		Short:             "Manage Gitea mirrors",
		SilenceUsage:      true,
		PersistentPreRunE: LoadConfig,
	}

	pFlags := cmd.PersistentFlags()

	pFlags.StringSliceP("config-file", "c", []string{"gitea-mirror.yaml"}, "configuration `file`s")
	pFlags.StringP("source.token", "S", "", "source API `token`")
	pFlags.StringP("target.token", "T", "", "target API `token`")
	pFlags.StringP("owner", "o", "", "default owner")
	pFlags.StringSliceP("labels", "l", nil, "filter repositories by label")

	cmd.AddCommand(
		cmdConfig(),
		cmdCreate(),
		cmdRecreate(),
		cmdDelete(),
		cmdStatus(),
		cmdSync(),
		cmdUpdate(),
	)

	return cmd
}

func LoadConfig(cmd *cobra.Command, args []string) (err error) {
	viper.BindPFlags(cmd.Flags())
	viper.SetEnvPrefix("GM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()      // read in environment variables that match bound variables
	viper.AllowEmptyEnv(true) // respect empty environment variables
	viper.BindEnv("source.token", "SOURCE_TOKEN")
	viper.BindEnv("target.token", "TARGET_TOKEN")

	config, err = cfg.LoadConfig(
		viper.GetStringSlice("config-file"),
	)
	if err != nil {
		return err
	}

	if labels := viper.GetStringSlice("labels"); len(labels) > 0 {
		config.Repositories = config.LabelledRepositories(labels)
	}

	return err
}

func newSource(config *cfg.Config) (server.Server, error) {
	switch config.Source.Type {
	case "github":
		return github.NewController(&config.Source), nil
	case "gitea":
		source, err := gitea.NewController(&config.Source)
		if err != nil {
			return nil, fmt.Errorf("NewController(%s): %w", config.Source.Url, err)
		}
		return source, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %q", config.Source.Type)
	}
}

// promptSourceToken prompts for a source token if any selected repository is private and none is configured.
func promptSourceToken(args []string) error {
	if config.Source.Token != "" {
		return nil
	}

	for repo := range config.FilteredRepositories(args) {
		if !*repo.PublicSource {
			if err := util.PromptForToken("Source API token", &config.Source.Token); err != nil {
				return fmt.Errorf("Source API token: %w", err)
			}
			return nil
		}
	}

	return nil
}
