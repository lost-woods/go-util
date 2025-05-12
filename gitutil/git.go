package gitutil

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/lost-woods/go-util/osutil"
)

var (
	log = osutil.GetLogger()
)

func CreateGit(config GitConfig) Git {
	auth := &http.BasicAuth{
		Username: config.Username,
		Password: config.Token,
	}

	git := Git{
		Config: config,
		Auth:   auth,
	}

	return git
}

func (git *Git) AddFile(path string, data []byte) {
	log.Infof("Adding file '%s'", path)
	osutil.WriteFile(osutil.JoinPath(git.Config.LocalPath, path), data)

	if _, err := git.WorkTree.Add(path); err != nil {
		log.Fatalf("Failed to stage file '%s': %v", path, err)
	}
}

func (git *Git) DeleteFile(path string) {
	log.Infof("Deleting file '%s'", path)
	if _, err := git.WorkTree.Remove(path); err != nil {
		log.Fatalf("Failed to stage deletion of file '%s': %v", path, err)
	}
}

func (git *Git) OpenOrClone() {
	if git.repoExists() {
		git.openRepository()
	} else {
		git.cloneRepository()
	}

	git.loadWorkTree()
}

func (git *Git) Checkout(branch string) {
	log.Infof("Checking out branch '%s'", branch)

	if git.tryCheckoutLocal(branch) {
		return
	}

	git.fetchRemoteBranches()

	if git.remoteBranchExists(branch) {
		git.checkoutRemoteBranch(branch)
		return
	}

	git.createAndCheckoutBranchFromHEAD(branch)
	git.setUpstreamTracking(branch)
}

func (git *Git) CommitAndPush(message string) string {
	status, err := git.WorkTree.Status()
	if err != nil {
		log.Fatalf("Failed to get worktree status: %v", err)
	}

	if status.IsClean() {
		log.Infof("No changes to commit")
		return ""
	} else {
		git.pullIfTracking()
	}

	hash := git.commitChanges(message)
	git.pushChanges()

	return hash
}
