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

func (c *Controller) GetOrg(r *config.Repository) *gitea.Organization {
	// Check cache first
	if org, cached := c.orgCache[r.Owner]; cached {
		return org
	}

	// Cache miss - query API
	org, _, err := c.client.GetOrg(r.Owner)
	if err != nil {
		// Cache the fact that org doesn't exist
		c.orgCache[r.Owner] = nil
		return nil
	}

	// Cache the successful result
	c.orgCache[r.Owner] = org
	return org
}

func (c *Controller) CreateOrg(orgName string, visibility gitea.VisibleType) (*gitea.Organization, error) {
	options := gitea.CreateOrgOption{
		Name:       orgName,
		FullName:   orgName,
		Visibility: visibility,
	}
	org, _, err := c.client.CreateOrg(options)
	if err != nil {
		return nil, err
	}
	return org, err
}

func (c *Controller) EnsureOrg(r *config.Repository) error {
	// Check if organization already exists (uses cache)
	org := c.GetOrg(r)
	if org != nil {
		return nil
	}

	// Determine visibility based on repository's PublicTarget setting
	visibility := gitea.VisibleTypePublic
	if !*r.PublicTarget {
		visibility = gitea.VisibleTypePrivate
	}

	// Create the organization with appropriate visibility
	org, err := c.CreateOrg(r.Owner, visibility)
	if err != nil {
		return fmt.Errorf("creating organization %s: %w", r.Owner, err)
	}

	// Update cache with newly created org
	c.orgCache[r.Owner] = org

	return nil
}

func (c *Controller) GetMirror(r *config.Repository) (*gitea.Repository, error) {
	repo, _, err := c.client.GetRepo(r.Owner, r.Name)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (c *Controller) CreateMirror(source server.Server, r *config.Repository) (*gitea.Repository, error) {
	// Ensure the organization exists before attempting to create the mirror
	if err := c.EnsureOrg(r); err != nil {
		return nil, fmt.Errorf("ensuring organization %s: %w", r.Owner, err)
	}

	cloneURL, err := source.GetCloneURL(r)
	if err != nil {
		return nil, err
	}

	var authToken string
	if !*r.PublicSource {
		authToken = source.GetToken()
	}

	options := gitea.MigrateRepoOption{
		RepoOwner:      r.Owner,
		RepoName:       r.Name,
		Private:        !*r.PublicTarget,
		CloneAddr:      cloneURL,
		AuthToken:      authToken,
		Mirror:         true,
		MirrorInterval: r.Interval.String(),
	}

	mirror, _, err := c.client.MigrateRepo(options)

	return mirror, err
}

func (c *Controller) UpdateMirror(r *config.Repository) (*gitea.Repository, error) {
	private := !*r.PublicTarget
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
