package wrapper

import (
	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
)

// RepositoriesList returns all repositories to which the user has access to
func (w *GitHubWrapper) GistsList(organization string) ([]*github.Gist, error) {
	logger := logrus.WithField("package", "wrapper").WithField("method", "GistsList")
	logger.Debugf(">> called with organization := [%s]", organization)
	defer logger.Debugf("<< done for organization := [%s]", organization)

	var (
		result = make([]*github.Gist, 0)
	)
	opt := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		gists, resp, err := w.client.Gists.List(w.ctx, organization, opt)
		if err != nil {
			return nil, err
		}
		result = append(result, gists...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return result, nil
}
