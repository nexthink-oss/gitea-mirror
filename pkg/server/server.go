package server

type Server interface {
	GetType() string
	GetToken() (token string)
	GetCloneURL(owner, name string) (string, error)
}
