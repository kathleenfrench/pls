package git

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/pkg/utils"
	"golang.org/x/oauth2"
)

// NewClient initializes a new github client
func NewClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

// NewEnterpriseClient creates a new github enterprise client
func NewEnterpriseClient(ctx context.Context, hostname string, token string) (*github.Client, error) {
	validBaseURL, err := utils.ValidateURL(hostname)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	tc := oauth2.NewClient(ctx, ts)
	uploadURL := fmt.Sprintf("https://%s/api/uploads/", hostname)
	client, err := github.NewEnterpriseClient(validBaseURL, uploadURL, tc)
	if err != nil {
		return nil, err
	}

	return client, nil
}
