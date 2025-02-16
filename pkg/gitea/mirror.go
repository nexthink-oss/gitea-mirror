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

func (c *Controller) GetMirror(r *config.Repository) (*gitea.Repository, error) {
	repo, _, err := c.client.GetRepo(r.Owner, r.Name)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (c *Controller) CreateMirror(source server.Server, r *config.Repository) (*gitea.Repository, error) {
	cloneURL, err := source.GetCloneURL(r)
	if err != nil {
		return nil, err
	}

	options := gitea.MigrateRepoOption{
		RepoOwner:      r.Owner,
		RepoName:       r.Name,
		Private:        !*r.Public,
		CloneAddr:      cloneURL,
		AuthToken:      source.GetToken(),
		Mirror:         true,
		MirrorInterval: r.Interval.String(),
	}

	mirror, _, err := c.client.MigrateRepo(options)

	return mirror, err
}

func (c *Controller) UpdateMirror(r *config.Repository) (*gitea.Repository, error) {
	private := !*r.Public
	interval := r.Interval.String()
	options := gitea.EditRepoOption{
		Private:        &private,
		MirrorInterval: &interval,
	}

	repo, _, err := c.client.EditRepo(r.Owner, r.Name, options)

	return repo, err
}

func (c *Controller) SyncMirror(r *config.Repository) error {
	_, err := c.client.MirrorSync(r.Owner, r.Name)

	return err
}

func (c *Controller) LastSynced(r *config.Repository) (*time.Time, error) {
	repo, _, err := c.client.GetRepo(r.Owner, r.Name)
	if err != nil {
		return nil, err
	}

	if !repo.Mirror {
		return nil, fmt.Errorf("Repository is not a mirror: %s/%s", r.Owner, r.Name)
	}

	return &repo.MirrorUpdated, nil
}

func (c *Controller) DeleteMirror(r *config.Repository) error {
	_, err := c.client.DeleteRepo(r.Owner, r.Name)

	return err
}
