package github

import (
	"context"

	"github.com/google/go-github/v68/github"
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

func (c Controller) GetCloneURL(owner, name string) (string, error) {
	repo, _, err := c.client.Repositories.Get(c.ctx, owner, name)

	return *repo.CloneURL, err
}
