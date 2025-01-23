package gitea

import (
	"context"

	"code.gitea.io/sdk/gitea"
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

func (g Controller) GetCloneURL(owner, name string) (string, error) {
	repo, _, err := g.client.GetRepo(owner, name)
	return repo.CloneURL, err
}
