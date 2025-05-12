package gitutil

import (
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	gogitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/lost-woods/go-util/osutil"
)

func (git *Git) repoExists() bool {
	return osutil.PathExists(osutil.JoinPath(git.Config.LocalPath, ".git"))
}

func (git *Git) openRepository() {
	log.Infof("Opening existing repository at '%s'", git.Config.LocalPath)
	repo, err := gogit.PlainOpen(git.Config.LocalPath)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}
	git.RepoObject = repo
}

func (git *Git) cloneRepository() {
	log.Infof("Cloning repository '%s'", git.Config.Repository)
	repo, err := gogit.PlainClone(git.Config.LocalPath, false, &gogit.CloneOptions{
		URL:  git.Config.Repository,
		Auth: git.Auth,
	})
	if err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	git.RepoObject = repo
}

func (git *Git) loadWorkTree() {
	workTree, err := git.RepoObject.Worktree()
	if err != nil || workTree == nil {
		log.Fatalf("Failed to get worktree: %v", err)
	}
	git.WorkTree = workTree
}

func (git *Git) tryCheckoutLocal(branch string) bool {
	err := git.WorkTree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: false,
	})
	if err == nil {
		log.Infof("Switched to existing local branch '%s'", branch)
		return true
	}
	return false
}

func (git *Git) fetchRemoteBranches() {
	err := git.RepoObject.Fetch(&gogit.FetchOptions{
		RemoteName: "origin",
		Auth:       git.Auth,
		Tags:       gogit.NoTags,
	})
	if err != nil && err != gogit.NoErrAlreadyUpToDate {
		log.Warnf("Fetch failed or not needed: %v", err)
	}
}

func (git *Git) remoteBranchExists(branch string) bool {
	_, err := git.RepoObject.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
	return err == nil
}

func (git *Git) checkoutRemoteBranch(branch string) {
	ref, err := git.RepoObject.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
	if err != nil {
		log.Fatalf("Remote branch '%s' not found: %v", branch, err)
	}

	localBranchRef := plumbing.NewBranchReferenceName(branch)
	err = git.WorkTree.Checkout(&gogit.CheckoutOptions{
		Hash:   ref.Hash(),
		Branch: localBranchRef,
		Create: true,
	})
	if err != nil {
		log.Fatalf("Failed to checkout remote branch '%s': %v", branch, err)
	}

	git.setUpstreamTracking(branch)

	log.Infof("Checked out remote branch '%s'", branch)
}

func (git *Git) createAndCheckoutBranchFromHEAD(branch string) {
	head, err := git.RepoObject.Head()
	if err != nil {
		log.Fatalf("Unable to retrieve HEAD: %v", err)
	}

	err = git.WorkTree.Checkout(&gogit.CheckoutOptions{
		Hash:   head.Hash(),
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: true,
	})
	if err != nil {
		log.Fatalf("Failed to create and checkout branch '%s': %v", branch, err)
	}
	log.Infof("Created new local branch '%s'", branch)
}

func (git *Git) setUpstreamTracking(branch string) {
	err := git.RepoObject.CreateBranch(&gogitconfig.Branch{
		Name:   branch,
		Remote: "origin",
		Merge:  plumbing.ReferenceName("refs/heads/" + branch),
	})
	if err != nil && err != gogit.ErrBranchExists {
		log.Fatalf("Failed to set tracking for branch '%s': %v", branch, err)
	}
	log.Infof("Set tracking for '%s' to origin/%s", branch, branch)
}

func (git *Git) pullIfTracking() {
	head, err := git.RepoObject.Head()
	if err != nil {
		log.Fatalf("Unable to get HEAD: %v", err)
	}
	branch := head.Name().Short()
	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)

	if _, err := git.RepoObject.Reference(remoteRef, true); err == nil {
		log.Infof("Pulling changes for branch '%s'", branch)
		err = git.WorkTree.Pull(&gogit.PullOptions{
			Auth:         git.Auth,
			RemoteName:   "origin",
			SingleBranch: true,
		})
		if err != nil && err != gogit.NoErrAlreadyUpToDate {
			log.Warnf("Pull failed: %v", err)
		}
	} else if err != plumbing.ErrReferenceNotFound {
		log.Fatalf("Failed checking remote ref: %v", err)
	}
}

func (git *Git) commitChanges(message string) string {
	commit, err := git.WorkTree.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  git.Config.Name,
			Email: git.Config.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Commit failed: %v", err)
	}

	commitObj, err := git.RepoObject.CommitObject(commit)
	if err != nil {
		log.Fatalf("Unable to get commit object: %v", err)
	}

	log.Infof("Created commit '%s'", commitObj.Hash.String())
	return commitObj.Hash.String()
}

func (git *Git) pushChanges() {
	head, _ := git.RepoObject.Head()
	branch := head.Name().Short()

	err := git.RepoObject.Push(&gogit.PushOptions{
		Auth: git.Auth,
		RefSpecs: []gogitconfig.RefSpec{
			gogitconfig.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)),
		},
	})
	if err != nil && err != gogit.NoErrAlreadyUpToDate {
		log.Fatalf("Push failed: %v", err)
	}

	log.Infof("Pushed branch '%s'", branch)
}
