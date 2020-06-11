## multi-gitter run

Run clones multiple repostories, run a script in that directory, and creates a PR with those changes.

### Synopsis

Run will clone down multiple repositories. For each of those repositories, the script will be run. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

```
multi-gitter run [script path] [flags]
```

### Options

```
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -h, --help                    help for run
  -R, --max-reviewers int       If this value is set, reviewers will be randomized
  -o, --org string              The name of the GitHub organization.
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.
```

### Options inherited from parent commands

```
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. (default "https://api.github.com/")
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories
