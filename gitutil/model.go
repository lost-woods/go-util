package gitutil

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Git struct {
	Config     GitConfig
	Auth       *http.BasicAuth
	WorkTree   *gogit.Worktree
	RepoObject *gogit.Repository
}

type GitConfig struct {
	Name       string
	Email      string
	Username   string
	Token      string
	Repository string
	LocalPath  string
}
