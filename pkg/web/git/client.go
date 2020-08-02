package git

import (
	"context"

	"github.com/google/go-github/v32/github"
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
