package gitea

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	ctx      context.Context
	client   *gitea.Client
	forge    config.Forge
	orgCache map[string]*gitea.Organization
}

func NewController(ctx context.Context, forge config.Forge) (*Controller, error) {
	client, err := gitea.NewClient(forge.GetUrl(), gitea.SetToken(forge.GetToken()))
	if err != nil {
		return nil, err
	}

	return &Controller{
		ctx:      ctx,
		client:   client,
		forge:    forge,
		orgCache: make(map[string]*gitea.Organization),
	}, nil
}

func (g *Controller) GetType() string {
	return "gitea"
}

func (g *Controller) GetToken() string {
	return g.forge.GetToken()
}

func (c *Controller) GetCloneURL(r *config.Repository) (string, error) {
	if c.forge != nil {
		if remoteUrl := c.forge.GetRemoteUrl(); remoteUrl != "" {
			return fmt.Sprintf("%s/%s/%s.git", remoteUrl, r.Owner, r.Name), nil
		}
	}

	repo, _, err := c.client.GetRepo(r.Owner, r.Name)
	if err != nil {
		return "", err
	}

	return repo.CloneURL, nil
}
