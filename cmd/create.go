package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

func cmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [<repository> ...]",
		Short: "Create Gitea mirrors",
		RunE:  CreateMirrors,
	}

	return cmd
}

func CreateMirrors(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()

	if config.Source.Token == "" {
		if err := util.PromptForToken("Source API token", &config.Source.Token); err != nil {
			return fmt.Errorf("Source API token: %w", err)
		}
	}

	if config.Target.Token == "" {
		if err := util.PromptForToken("Target API token", &config.Target.Token); err != nil {
			return fmt.Errorf("Target API token: %w", err)
		}
	}

	source, err := newSource(config)
	if err != nil {
		return err
	}

	target, err := gitea.NewController(&config.Target)
	if err != nil {
		return fmt.Errorf("NewController(%s): %w", config.Target.Url, err)
	}

	for repo := range config.FilteredRepositories(args) {
		if _, err = target.CreateMirror(ctx, source, &repo); err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success())
		}
	}

	return nil
}
