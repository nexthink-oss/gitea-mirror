package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

func cmdRecreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recreate [<repository> ...]",
		Short: "Recreate Gitea mirrors",
		RunE:  RecreateMirrors,
	}

	return cmd
}

func RecreateMirrors(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()

	if err := promptSourceToken(args); err != nil {
		return err
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
		if err = target.DeleteMirror(ctx, &repo); err != nil {
			fmt.Println(repo.Failure(fmt.Errorf("deleting: %w", err)))
			continue
		}
		if _, err = target.CreateMirror(ctx, source, &repo); err != nil {
			fmt.Println(repo.Failure(fmt.Errorf("creating: %w", err)))
		} else {
			fmt.Println(repo.Success())
		}
	}

	return nil
}
