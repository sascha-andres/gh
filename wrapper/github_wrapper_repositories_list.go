package wrapper

import (
	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
)

// RepositoriesList returns all repositories to which the user has access to
func (w *GitHubWrapper) RepositoriesList(affiliation, visibility string) ([]*github.Repository, error) {
	logger := logrus.WithField("package", "wrapper").WithField("method", "RepositoriesList")
	logger.Debugf(">> called with affiliation := [%s] visibility := [%s]", affiliation, visibility)
	defer logger.Debugf("<< done for affiliation := [%s] visibility := [%s]", affiliation, visibility)

	var (
		result = make([]*github.Repository, 0)
	)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Affiliation: affiliation,
		Visibility:  visibility,
	}

	for {
		repos, resp, err := w.client.Repositories.List(w.ctx, "", opt)
		if err != nil {
			return nil, err
		}
		result = append(result, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return result, nil
}
