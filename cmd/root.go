package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/nexthink-oss/gitea-mirror/pkg/config"
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

	pFlags.StringArrayP("config-path", "P", []string{"."}, "configuration file path")
	pFlags.StringP("config-name", "C", "gitea-mirror", "configuration file name (without extension)")
	pFlags.StringP("source.token", "S", "", "source API token")
	pFlags.StringP("target.token", "T", "", "target API token")
	pFlags.StringP("owner", "o", "", "default owner")

	cmd.AddCommand(
		cmdConfig(),
		cmdCreate(),
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
		viper.GetString("config-name"),
		viper.GetStringSlice("config-path")...,
	)
	if err != nil {
		return err
	}

	return err
}
