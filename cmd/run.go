package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// RunCmd is the main command that runs a script for multiple repositories and creates PRs with the changes made
var RunCmd = &cobra.Command{
	Use:   "run [script path]",
	Short: "Run clones multiple repostories, run a script in that directory, and creates a PR with those changes.",
	Long:  "Run will clone down multiple repositories. For each of those repositories, the script will be run. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func init() {
	RunCmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	RunCmd.Flags().StringP("org", "o", "", "The name of the GitHub organization.")
	RunCmd.Flags().StringP("pr-title", "t", "", "The title of the PR. Will default to the first line of the commit message if none is set.")
	RunCmd.Flags().StringP("pr-body", "b", "", "The body of the commit message. Will default to everything but the first line of the commit message if none is set.")
	RunCmd.Flags().StringP("commit-message", "m", "", "The commit message. Will default to title + body if none is set.")
	RunCmd.Flags().StringSliceP("reviewers", "r", nil, "The username of the reviewers to be added on the pull request.")
	RunCmd.Flags().IntP("max-reviewers", "R", 0, "If this value is set, reviewers will be randomized")
}

func run(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	ghBaseUrl, _ := flag.GetString("gh-base-url")
	token, _ := flag.GetString("token")
	branchName, _ := flag.GetString("branch")
	org, _ := flag.GetString("org")
	prTitle, _ := flag.GetString("pr-title")
	prBody, _ := flag.GetString("pr-body")
	commitMessage, _ := flag.GetString("commit-message")
	reviewers, _ := flag.GetStringSlice("reviewers")
	maxReviewers, _ := flag.GetInt("max-reviewers")

	programPath := flag.Arg(0)

	if token != "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		fmt.Println("Either the --token flag or the GITHUB_TOKEN environment variable has to be set.")
		cmd.Usage()
		os.Exit(1)
	}

	if org == "" {
		fmt.Println("No organisation set.")
		cmd.Usage()
		os.Exit(1)
	}

	// Set commit message based on pr title and body or the reverse
	if commitMessage == "" && prTitle == "" {
		fmt.Println("Pull request title or commit message must be set.")
		cmd.Usage()
		os.Exit(1)
	} else if commitMessage == "" {
		commitMessage = prTitle
		if prBody != "" {
			commitMessage += "\n" + prBody
		}
	} else if prTitle == "" {
		split := strings.SplitN(commitMessage, "\n", 2)
		prTitle = split[0]
		if prBody == "" && len(split) == 2 {
			prBody = split[2]
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(workingDir)
	}

	vc, err := github.New(token, ghBaseUrl)
	if err != nil {
		log.Fatalln(workingDir)
	}

	runner := multigitter.Runner{
		ScriptPath:    path.Join(workingDir, programPath),
		FeatureBranch: branchName,

		VersionController: vc,

		OrgName:          org,
		CommitMessage:    commitMessage,
		PullRequestTitle: prTitle,
		PullRequestBody:  prBody,
		Reviewers:        reviewers,
		MaxReviewers:     maxReviewers,
	}

	err = runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}