package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

var updateCmd = &cobra.Command{
	Use:   "update [<repository> ...]",
	Short: "Update Gitea mirrors",
	RunE:  UpdateMirrors,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func UpdateMirrors(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()

	if config.Target.Token == "" {
		if err := util.PromptForToken("Target API token", &config.Target.Token); err != nil {
			return fmt.Errorf("Target API token: %w", err)
		}
	}

	target, err := gitea.NewController(ctx, config.Target.Url, config.Target.Token)
	if err != nil {
		return fmt.Errorf("NewController(%s): %w", config.Target.Url, err)
	}

	for repo := range config.FilteredRepositories(args) {
		if _, err = target.UpdateMirror(&repo); err != nil {
			fmt.Println(repo.Failure(err))
		} else {
			fmt.Println(repo.Success())
		}
	}

	return nil
}
