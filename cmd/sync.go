package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

func cmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [<repository> ...]",
		Short: "Sync Gitea mirrors",
		RunE:  SyncMirrors,
	}

	return cmd
}

func SyncMirrors(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()

	if config.Target.Token == "" {
		if err := util.PromptForToken("Target API token", &config.Target.Token); err != nil {
			return err
		}
	}

	target, err := gitea.NewController(ctx, &config.Target)
	if err != nil {
		return err
	}

	for repo := range config.FilteredRepositories(args) {
		if err := target.SyncMirror(&repo); err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success())
		}
	}

	return nil
}
