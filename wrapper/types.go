package wrapper

import (
	"context"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type (
	// GitHubWrapper is an anchor to GitHub wrapper functions
	GitHubWrapper struct {
		client *github.Client

		ctx context.Context
	}
)

// NewGitHubWrapper returns a new git hub wrapper object
func NewGitHubWrapper(token string) (*GitHubWrapper, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubWrapper{
		client: github.NewClient(tc),
		ctx:    ctx,
	}, nil
}
