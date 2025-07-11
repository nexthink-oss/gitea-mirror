package cmd

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
	out, err := yaml.MarshalWithOptions(config, yaml.IndentSequence(true))
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
