package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the resolved configuration",
	RunE:  ShowConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func ShowConfig(cmd *cobra.Command, args []string) error {
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
