package multigitter

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git"
)

// RepoGetter fetches repositories
type RepoGetter interface {
	GetRepositories() ([]domain.Repository, error)
}

// PullRequestCreator creates pull requests
type PullRequestCreator interface {
	CreatePullRequest(repo domain.Repository, newPR domain.NewPullRequest) error
}

// Runner conains fields to be able to do the run
type Runner struct {
	ScriptPath    string // Must be absolute path
	FeatureBranch string

	RepoGetter         RepoGetter
	PullRequestCreator PullRequestCreator

	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run() error {
	repos, err := r.RepoGetter.GetRepositories()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		log.Printf("Cloning and running script on: %s\n", repo.GetURL())
		err := r.runSingleRepo(repo.GetURL())
		switch {
		case err == domain.ExitCodeError:
			log.Printf("Got exit code when running %s\n", repo.GetURL())
			continue
		case err == domain.NoChangeError:
			log.Printf("No change done on the repo by the script when running: %s\n", repo.GetURL())
			continue
		case err != nil:
			return err
		}

		err = r.PullRequestCreator.CreatePullRequest(repo, domain.NewPullRequest{
			Title:     r.PullRequestTitle,
			Body:      r.PullRequestBody,
			Head:      r.FeatureBranch,
			Base:      repo.GetBranch(),
			Reviewers: r.Reviewers,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Runner) runSingleRepo(url string) error {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := git.Git{
		Directory: tmpDir,
		Repo:      url,
		NewBranch: r.FeatureBranch,
	}

	err = sourceController.Clone()
	if err != nil {
		return err
	}

	// Run the command that might or might not change the content of the repo
	// If the command return a non zero exit code, abort.
	cmd := exec.Command(r.ScriptPath)
	cmd.Dir = tmpDir

	reader, writer := io.Pipe()
	cmd.Stdout = writer
	cmd.Stderr = writer

	// Print each line that is outputted by the script
	go func() {
		buf := bufio.NewReader(reader)
		for {
			line, err := buf.ReadString('\n')
			if line != "" {
				log.Printf("Script output: %s", line)
			}
			if err != nil {
				return
			}
		}
	}()

	err = cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return domain.ExitCodeError
		}
		return err
	}

	err = sourceController.Commit(r.CommitMessage)
	if err != nil {
		return err
	}

	return nil
}
