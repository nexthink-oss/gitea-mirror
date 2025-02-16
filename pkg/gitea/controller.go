package gitea

import (
	"context"

	"code.gitea.io/sdk/gitea"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	ctx    context.Context
	client *gitea.Client
	token  string
}

func NewController(ctx context.Context, url, token string) (*Controller, error) {
	client, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		return nil, err
	}

	return &Controller{
		ctx:    ctx,
		client: client,
		token:  token,
	}, nil
}

func (g Controller) GetType() string {
	return "gitea"
}

func (g Controller) GetToken() string {
	return g.token
}

func (c Controller) GetCloneURL(r *config.Repository) (string, error) {
	repo, _, err := c.client.GetRepo(r.Owner, r.Name)
	return repo.CloneURL, err
}
