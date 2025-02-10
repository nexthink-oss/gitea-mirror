package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Gitea mirrors",
	RunE:  SyncMirrors,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func SyncMirrors(cmd *cobra.Command, args []string) (err error) {
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
		if err := target.SyncMirror(&repo); err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success())
		}
	}

	return nil
}
