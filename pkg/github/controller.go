package github

import (
	"context"

	"github.com/google/go-github/v68/github"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	ctx    context.Context
	client *github.Client
	token  string
}

func NewController(ctx context.Context, token string) *Controller {
	client := github.NewClient(nil).WithAuthToken(token)

	return &Controller{
		ctx:    ctx,
		client: client,
		token:  token,
	}
}

func (c Controller) GetType() string {
	return "github"
}

func (c Controller) GetToken() string {
	return c.token
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

	return *repo.CloneURL, err
}
