package gitea

import (
	"context"
	"fmt"

	"gitea.dev/sdk"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Controller struct {
	client     *gitea.Client
	forge      config.Forge
	orgCache   map[string]*gitea.Organization
	mirrorAuth *bool
}

func NewController(forge config.Forge) (*Controller, error) {
	client, err := gitea.NewClient(forge.GetUrl(), gitea.SetToken(forge.GetToken()))
	if err != nil {
		return nil, err
	}

	return &Controller{
		client:   client,
		forge:    forge,
		orgCache: make(map[string]*gitea.Organization),
	}, nil
}

func (c *Controller) GetType() string {
	return "gitea"
}

func (c *Controller) GetToken() string {
	return c.forge.GetToken()
}

func (c *Controller) GetCloneURL(ctx context.Context, r *config.Repository) (string, error) {
	if c.forge != nil {
		if remoteUrl := c.forge.GetRemoteUrl(); remoteUrl != "" {
			return fmt.Sprintf("%s/%s/%s.git", remoteUrl, r.Owner, r.Name), nil
		}
	}

	repo, _, err := c.client.Repositories.GetRepo(ctx, r.Owner, r.Name)
	if err != nil {
		return "", err
	}

	return repo.CloneURL, nil
}

// supportsMirrorAuth reports whether the server accepts mirror credentials on repository edit (Gitea >= 1.27).
func (c *Controller) supportsMirrorAuth(ctx context.Context) bool {
	if c.mirrorAuth == nil {
		supported := c.client.CheckServerVersionConstraint(ctx, ">=1.27.0") == nil
		c.mirrorAuth = &supported
	}

	return *c.mirrorAuth
}
