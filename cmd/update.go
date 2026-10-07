package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nexthink-oss/gitea-mirror/pkg/gitea"
	"github.com/nexthink-oss/gitea-mirror/pkg/server"
	"github.com/nexthink-oss/gitea-mirror/pkg/util"
)

func cmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [<repository> ...]",
		Short: "Update Gitea mirrors",
		RunE:  UpdateMirrors,
	}

	cmd.Flags().Bool("skip-credentials", false, "do not update mirror source credentials")

	return cmd
}

func UpdateMirrors(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	var source server.Server

	if !viper.GetBool("skip-credentials") {
		if err := promptSourceToken(args); err != nil {
			return err
		}

		source, err = newSource(config)
		if err != nil {
			return err
		}
	}

	if config.Target.Token == "" {
		if err := util.PromptForToken("Target API token", &config.Target.Token); err != nil {
			return fmt.Errorf("Target API token: %w", err)
		}
	}

	target, err := gitea.NewController(&config.Target)
	if err != nil {
		return fmt.Errorf("NewController(%s): %w", config.Target.Url, err)
	}

	warned := false
	for repo := range config.FilteredRepositories(args) {
		_, err = target.UpdateMirror(ctx, source, &repo)
		var notUpdated *gitea.TokenNotUpdated
		switch {
		case errors.As(err, &notUpdated):
			fmt.Println(repo.Success())
			if !warned {
				fmt.Println("warning: Gitea < 1.27 cannot update mirror credentials; use `recreate` to rotate the source token")
				warned = true
			}
		case err != nil:
			fmt.Println(repo.Failure(err))
		default:
			fmt.Println(repo.Success())
		}
	}

	return nil
}
