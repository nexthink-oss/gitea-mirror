package cmd

import (
	"fmt"
	"time"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of the mirrors",
	RunE:  Status,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func Status(cmd *cobra.Command, args []string) error {
	var ctx = cmd.Context()

	if config.Target.Token == "" {
		if err := util.PromptForToken("Server API token", &config.Target.Token); err != nil {
			return err
		}
	}

	target, err := gitea.NewController(ctx, config.Target.Url, config.Target.Token)
	if err != nil {
		return err
	}

	for _, repo := range config.Repositories {
		updated, err := target.StatusMirror(repo)
		if err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success(updated.UTC().Format(time.RFC3339)))
		}
	}

	return nil
}
