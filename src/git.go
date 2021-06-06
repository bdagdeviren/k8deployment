package src

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
	"os"
)

func GetRemoteCommitId(url,branch,token string) string {
	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	commitId := ""

	for _, ref := range refs {
		if ref.Name().IsBranch() {
			if ref.Name().Short() == branch {
				commitId = ref.Hash().String()
			}
		}
	}

	return commitId
}

func CloneGitRepository(url,branch,token string) string {
	r, err := git.PlainClone("deployment", false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "abc123",
			Password: token,
		},
		URL: url,
		RemoteName: branch,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil{
		log.Fatalln(err.Error())
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil{
		log.Fatalln(err.Error())
	}
	// ... retrieving the commit object
	commit := ref.Hash().String()

	return commit
}

func GetLocalCommitId() string {
	r, err := git.PlainOpen("deployment")
	if err != nil{
		log.Fatalln(err.Error())
	}

	ref, err := r.Head()
	if err != nil{
		log.Fatalln(err.Error())
	}

	return ref.Hash().String()
}

func PullGitRepository() string {
	r, err := git.PlainOpen("deployment")
	if err != nil{
		log.Fatalln(err.Error())
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil{
		log.Fatalln(err.Error())
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil{
		log.Fatalln(err.Error())
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil{
		log.Fatalln(err.Error())
	}

	commit := ref.Hash().String()
	return commit
}

func CloneOrPullRepository(url,branch,token string) bool {
	deploy := false

	if _, err := os.Stat("deployment"); os.IsNotExist(err) {
		log.Println("Cloning Repository!")
		cloneCommit := CloneGitRepository(url, branch, token)
		log.Printf("Cloned Repository - CommitId: %s\n",cloneCommit)
		deploy = true
	}else {
		localCommitId := GetLocalCommitId()
		remoteCommitId := GetRemoteCommitId(url,branch,token)

		log.Printf("Local CommitId: %s\n",localCommitId)
		log.Printf("Remote CommitId: %s\n",remoteCommitId)

		if localCommitId != remoteCommitId{
			log.Println("Pulling Repository!")
			gettingCommitId := PullGitRepository()
			log.Printf("Pulled Repository - CommitId: %s\n",gettingCommitId)
			deploy = true
		}
	}

	return deploy
}
