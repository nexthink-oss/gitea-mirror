package server

import (
	"context"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type Server interface {
	GetType() string
	GetToken() (token string)
	GetCloneURL(ctx context.Context, r *config.Repository) (string, error)
}
