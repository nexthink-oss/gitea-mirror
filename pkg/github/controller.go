package github

import (
	"context"

	"github.com/google/go-github/v71/github"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	ctx    context.Context
	client *github.Client
	forge  config.Forge
}

func NewController(ctx context.Context, forge config.Forge) *Controller {
	client := github.NewClient(nil).WithAuthToken(forge.GetToken())

	return &Controller{
		ctx:    ctx,
		client: client,
		forge:  forge,
	}
}

func (c Controller) GetType() string {
	return "github"
}

func (c Controller) GetToken() string {
	return c.forge.GetToken()
}

func (c Controller) IsPrivate(r *config.Repository) (bool, error) {
	repo, _, err := c.client.Repositories.Get(c.ctx, r.Owner, r.Name)
	if err != nil {
		return true, err
	}

	return *repo.Private, nil
}

func (c Controller) GetCloneURL(r *config.Repository) (string, error) {
	repo, _, err := c.client.Repositories.Get(c.ctx, r.Owner, r.Name)
	if err != nil {
		return "", err
	} else {
		return *repo.CloneURL, nil
	}
}
