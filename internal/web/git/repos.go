package gitpls

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kathleenfrench/pls/pkg/web/git"
)

// FetchReposInOrganization fetches repositories in an organization
func FetchReposInOrganization(organization string, settings config.Settings, useEnterprise bool) ([]*github.Repository, error) {
	var allOrgRepos []*github.Repository
	var gc *github.Client
	var err error

	ctx := context.Background()

	if useEnterprise {
		gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
		if err != nil {
			return nil, err
		}
	} else {
		gc = git.NewClient(ctx, settings.GitToken)
	}

	opts := github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := gc.Repositories.ListByOrg(ctx, organization, &opts)
		if err != nil {
			return nil, err
		}

		allOrgRepos = append(allOrgRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allOrgRepos, nil
}

// FetchUserRepos fetches repositories by user
func FetchUserRepos(username string, settings config.Settings, useEnterprise bool) ([]*github.Repository, error) {
	var gc *github.Client
	var err error

	ctx := context.Background()

	if useEnterprise {
		gc, err = git.NewEnterpriseClient(ctx, settings.GitEnterpriseHostname, settings.GitEnterpriseToken)
		if err != nil {
			return nil, err
		}
	} else {
		gc = git.NewClient(ctx, settings.GitToken)
	}

	opts := &github.RepositoryListOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allRepos []*github.Repository

	for {
		repos, resp, err := gc.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	if username == "" {
		username = "you"
	}

	gui.Log(":eyes:", fmt.Sprintf("%d repositories found for %s", len(allRepos), username), nil)
	return allRepos, nil
}
