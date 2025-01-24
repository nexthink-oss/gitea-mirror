package server

import "github.com/nexthink-oss/gitea-mirror/pkg/config"

type Server interface {
	GetType() string
	GetToken() (token string)
	GetCloneURL(*config.Repository) (string, error)
}
