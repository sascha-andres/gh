package wrapper

import (
	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
)

// RepositoriesListByOrganization returns all repositories inside an organization
func (w *GitHubWrapper) RepositoriesListByOrganization(organization string) ([]*github.Repository, error) {
	logger := logrus.WithField("package", "wrapper").WithField("method", "RepositoriesListByOrganization")
	logger.Debugf(">> called with organization := [%s]", organization)
	logger.Debugf("<< done for organization := [%s]", organization)

	var (
		result = make([]*github.Repository, 0)
	)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		repos, resp, err := w.client.Repositories.List(w.ctx, "", opt)
		if err != nil {
			return nil, err
		}
		for _, r := range repos {
			if *r.Owner.Login == organization {
				result = append(result, r)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return result, nil
}
