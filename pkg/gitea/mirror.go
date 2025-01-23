package gitea

import (
	"fmt"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
	"github.com/nexthink-oss/gitea-mirror/pkg/server"
)

type RepositoryNotMirror struct{}

func (e *RepositoryNotMirror) Error() string {
	return "Repository is not a mirror"
}

func (c *Controller) GetMirror(m config.Repository) (*gitea.Repository, error) {
	repo, _, err := c.client.GetRepo(m.Owner, m.Name)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (c *Controller) CreateMirror(source server.Server, r config.Repository) (*gitea.Repository, error) {
	cloneURL, err := source.GetCloneURL(r.Owner, r.Name)
	if err != nil {
		return nil, err
	}

	options := gitea.MigrateRepoOption{
		RepoOwner:      r.Owner,
		RepoName:       r.Name,
		Private:        r.Private,
		CloneAddr:      cloneURL,
		AuthToken:      source.GetToken(),
		Mirror:         true,
		MirrorInterval: "0",
	}

	mirror, _, err := c.client.MigrateRepo(options)

	return mirror, err
}

func (c *Controller) UpdateMirror(m config.Repository) error {
	_, err := c.client.MirrorSync(m.Owner, m.Name)

	return err
}

func (c *Controller) StatusMirror(m config.Repository) (*time.Time, error) {
	repo, _, err := c.client.GetRepo(m.Owner, m.Name)
	if err != nil {
		return nil, err
	}

	if !repo.Mirror {
		return nil, fmt.Errorf("Repository is not a mirror: %s/%s", m.Owner, m.Name)
	}

	return &repo.MirrorUpdated, nil
}
