package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func cmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Print the resolved configuration",
		RunE:  ShowConfig,
	}

	return cmd
}

func ShowConfig(cmd *cobra.Command, args []string) error {
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
