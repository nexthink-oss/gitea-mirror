package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Print the configuration",
	RunE:  DumpConfig,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func DumpConfig(cmd *cobra.Command, args []string) error {
	// config, err := config.LoadConfig()
	// if err != nil {
	// 	return err
	// }

	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
