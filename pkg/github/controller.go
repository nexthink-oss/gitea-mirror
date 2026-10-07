package github

import (
	"context"

	"github.com/google/go-github/v72/github"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	client *github.Client
	forge  config.Forge
}

func NewController(forge config.Forge) *Controller {
	client := github.NewClient(nil)
	if token := forge.GetToken(); token != "" {
		client = client.WithAuthToken(token)
	}

	return &Controller{
		client: client,
		forge:  forge,
	}
}

func (c *Controller) GetType() string {
	return "github"
}

func (c *Controller) GetToken() string {
	return c.forge.GetToken()
}

func (c *Controller) IsPrivate(ctx context.Context, r *config.Repository) (bool, error) {
	repo, _, err := c.client.Repositories.Get(ctx, r.Owner, r.Name)
	if err != nil {
		return true, err
	}

	return *repo.Private, nil
}

func (c *Controller) GetCloneURL(ctx context.Context, r *config.Repository) (string, error) {
	repo, _, err := c.client.Repositories.Get(ctx, r.Owner, r.Name)
	if err != nil {
		return "", err
	}

	return *repo.CloneURL, nil
}
