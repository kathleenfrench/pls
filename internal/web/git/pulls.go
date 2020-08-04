package gitpls

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/web/git"
)

// PullRequestSearchFlags are optional flags that can be set on a PR search
type PullRequestSearchFlags struct {
	Closed   bool
	Language string
	Draft    bool
}

// FetchUserPullRequestsEverywhere search all of github for user's pull requests
func FetchUserPullRequestsEverywhere(settings config.Settings) ([]*github.Issue, error) {
	var allPRs []*github.Issue

	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	ctx := context.Background()
	gc := git.NewClient(ctx, settings.GitToken)
	query := fmt.Sprintf("author:%s type:pr", settings.GitUsername)

	for {
		prs, resp, err := gc.Search.Issues(ctx, query, &opts)
		if err != nil {
			return nil, err
		}

		allPRs = append(allPRs, prs.Issues...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allPRs, nil
}
