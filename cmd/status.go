package cmd

import (
	"fmt"
	"time"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
	"github.com/spf13/cobra"
)

func cmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [<repository> ...]",
		Short: "Print the status of the mirrors",
		RunE:  Status,
	}

	return cmd
}

func Status(cmd *cobra.Command, args []string) error {
	var ctx = cmd.Context()

	if config.Target.Token == "" {
		if err := util.PromptForToken("Target API token", &config.Target.Token); err != nil {
			return err
		}
	}

	target, err := gitea.NewController(ctx, config.Target.Url, config.Target.Token)
	if err != nil {
		return err
	}

	for repo := range config.FilteredRepositories(args) {
		synced, err := target.LastSynced(&repo)
		if err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success(synced.UTC().Format(time.RFC3339)))
		}
	}

	return nil
}
