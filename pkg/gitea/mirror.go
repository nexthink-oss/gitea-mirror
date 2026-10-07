package gitea

import (
	"context"
	"fmt"
	"time"

	"gitea.dev/sdk"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
	"github.com/nexthink-oss/gitea-mirror/pkg/server"
)

// mirrorUsername matches the username Gitea stores for token-authenticated migrations.
const mirrorUsername = "oauth2"

type RepositoryNotMirror struct{}

func (e *RepositoryNotMirror) Error() string {
	return "Repository is not a mirror"
}

type TokenNotUpdated struct{}

func (e *TokenNotUpdated) Error() string {
	return "Mirror token not updated: requires Gitea >= 1.27"
}

func (c *Controller) GetOrg(ctx context.Context, r *config.Repository) *gitea.Organization {
	// Check cache first
	if org, cached := c.orgCache[r.Owner]; cached {
		return org
	}

	// Cache miss - query API
	org, _, err := c.client.Organizations.GetOrg(ctx, r.Owner)
	if err != nil {
		// Cache the fact that org doesn't exist
		c.orgCache[r.Owner] = nil
		return nil
	}

	// Cache the successful result
	c.orgCache[r.Owner] = org
	return org
}

func (c *Controller) CreateOrg(ctx context.Context, orgName string, visibility gitea.VisibleType) (*gitea.Organization, error) {
	options := gitea.CreateOrgOption{
		Name:       orgName,
		FullName:   orgName,
		Visibility: visibility,
	}
	org, _, err := c.client.Organizations.CreateOrg(ctx, options)
	if err != nil {
		return nil, err
	}
	return org, err
}

func (c *Controller) EnsureOrg(ctx context.Context, r *config.Repository) error {
	// Check if organization already exists (uses cache)
	org := c.GetOrg(ctx, r)
	if org != nil {
		return nil
	}

	// Determine visibility based on repository's PublicTarget setting
	visibility := gitea.VisibleTypePublic
	if !*r.PublicTarget {
		visibility = gitea.VisibleTypePrivate
	}

	// Create the organization with appropriate visibility
	org, err := c.CreateOrg(ctx, r.Owner, visibility)
	if err != nil {
		return fmt.Errorf("creating organization %s: %w", r.Owner, err)
	}

	// Update cache with newly created org
	c.orgCache[r.Owner] = org

	return nil
}

func (c *Controller) GetMirror(ctx context.Context, r *config.Repository) (*gitea.Repository, error) {
	repo, _, err := c.client.Repositories.GetRepo(ctx, r.Owner, r.Name)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (c *Controller) CreateMirror(ctx context.Context, source server.Server, r *config.Repository) (*gitea.Repository, error) {
	// Ensure the organization exists before attempting to create the mirror
	if err := c.EnsureOrg(ctx, r); err != nil {
		return nil, fmt.Errorf("ensuring organization %s: %w", r.Owner, err)
	}

	cloneURL, err := source.GetCloneURL(ctx, r)
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

	mirror, _, err := c.client.Repositories.MigrateRepo(ctx, options)

	return mirror, err
}

// UpdateMirror applies the configured visibility and interval and, unless source is nil or public,
// rotates the mirror's source token.
func (c *Controller) UpdateMirror(ctx context.Context, source server.Server, r *config.Repository) (*gitea.Repository, error) {
	private := !*r.PublicTarget
	interval := r.Interval.String()
	options := gitea.EditRepoOption{
		Private:        &private,
		MirrorInterval: &interval,
	}

	if source == nil || *r.PublicSource {
		repo, _, err := c.client.Repositories.EditRepo(ctx, r.Owner, r.Name, options)
		return repo, err
	}

	if !c.supportsMirrorAuth(ctx) {
		repo, _, err := c.client.Repositories.EditRepo(ctx, r.Owner, r.Name, options)
		if err != nil {
			return repo, err
		}
		return repo, &TokenNotUpdated{}
	}

	username := mirrorUsername
	token := source.GetToken()

	return c.editMirror(ctx, r.Owner, r.Name, editMirrorOption{
		EditRepoOption: options,
		MirrorUsername: &username,
		MirrorToken:    &token,
	})
}

func (c *Controller) SyncMirror(ctx context.Context, r *config.Repository) error {
	_, err := c.client.Repositories.MirrorSync(ctx, r.Owner, r.Name)

	return err
}

func (c *Controller) LastSynced(ctx context.Context, r *config.Repository) (*time.Time, error) {
	repo, _, err := c.client.Repositories.GetRepo(ctx, r.Owner, r.Name)
	if err != nil {
		return nil, err
	}

	if !repo.Mirror {
		return nil, fmt.Errorf("Repository is not a mirror: %s/%s", r.Owner, r.Name)
	}

	return &repo.MirrorUpdated, nil
}

func (c *Controller) DeleteMirror(ctx context.Context, r *config.Repository) error {
	_, err := c.client.Repositories.DeleteRepo(ctx, r.Owner, r.Name)

	return err
}
