package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/nexthink-oss/gitea-mirror/pkg/config"
)

var config *cfg.Config

var (
	version string = "snapshot"
	commit  string = "unknown"
	date    string = "unknown"
)

var rootCmd = &cobra.Command{
	Use:               "gitea-mirror",
	Short:             "Manage Gitea mirrors",
	SilenceUsage:      true,
	PersistentPreRunE: LoadConfig,
	Version:           fmt.Sprintf("%s-%s (built %s)", version, commit, date),
}

func init() {
	cobra.OnInitialize(initViper)

	pFlags := rootCmd.PersistentFlags()

	pFlags.StringArrayP("config-path", "P", []string{"."}, "configuration file path")
	viper.BindPFlag("config-path", pFlags.Lookup("config-path"))

	pFlags.StringP("config-name", "C", "gitea-mirror", "configuration file name")
	viper.BindPFlag("config-name", pFlags.Lookup("config-name"))

	pFlags.StringP("source-token", "S", "", "source API token")
	viper.BindPFlag("source.token", pFlags.Lookup("source-token"))
	viper.BindEnv("source.token", "SOURCE_TOKEN")

	pFlags.StringP("target-token", "T", "", "target API token")
	viper.BindPFlag("target.token", pFlags.Lookup("target-token"))
	viper.BindEnv("target.token", "TARGET_TOKEN")

	pFlags.StringP("owner", "o", "", "default owner")
	viper.BindPFlag("owner", pFlags.Lookup("owner"))
}

func initViper() {
	viper.SetEnvPrefix("GM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()      // read in environment variables that match bound variables
	viper.AllowEmptyEnv(true) // respect empty environment variables
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func LoadConfig(cmd *cobra.Command, args []string) (err error) {
	config, err = cfg.LoadConfig(
		viper.GetString("config-name"),
		viper.GetStringSlice("config-path")...,
	)
	if err != nil {
		return err
	}

	return err
}
