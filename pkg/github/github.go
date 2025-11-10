package github

import (
	"context"

	gh "github.com/google/go-github/v78/github"
)

const maxPerPage = 100

type Client interface {
	// Migrations
	GetMigrationArchiveURL(ctx context.Context, organization string, organizationID int64) (string, error)
	GetMigrationStatus(ctx context.Context, organization string, migrationID int64) (*gh.Migration, error)
	StartMigration(ctx context.Context, organization string, repoNames []string) (*gh.Migration, error)

	// Repositories
	ListOrgRepos(ctx context.Context, organization string, visibility string) ([]*gh.Repository, error)
}

type defaultClient struct {
	githubClient *gh.Client
}

func NewClient(gitHubClient *gh.Client) Client {
	return &defaultClient{
		githubClient: gitHubClient,
	}
}

func (c *defaultClient) GetMigrationArchiveURL(ctx context.Context, organization string, organizationID int64) (string, error) {
	return c.githubClient.Migrations.MigrationArchiveURL(ctx, organization, organizationID)
}

func (c *defaultClient) GetMigrationStatus(ctx context.Context, organization string, migrationID int64) (*gh.Migration, error) {
	migration, _, err := c.githubClient.Migrations.MigrationStatus(ctx, organization, migrationID)
	if err != nil {
		return nil, err
	}

	return migration, nil
}

func (c *defaultClient) StartMigration(ctx context.Context, organization string, repoNames []string) (*gh.Migration, error) {
	migration, _, err := c.githubClient.Migrations.StartMigration(ctx, organization, repoNames, &gh.MigrationOptions{
		ExcludeAttachments: true,
		ExcludeReleases:    true,
		Exclude:            []string{"repositories"},
	})
	if err != nil {
		return nil, err
	}

	return migration, nil
}

func (c *defaultClient) ListOrgRepos(ctx context.Context, organization string, visibility string) ([]*gh.Repository, error) {
	opts := &gh.RepositoryListByOrgOptions{
		Type: visibility,
		ListOptions: gh.ListOptions{
			PerPage: maxPerPage,
			Page:    1,
		},
	}

	var allRepos []*gh.Repository
	for {
		repos, resp, err := c.githubClient.Repositories.ListByOrg(ctx, organization, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allRepos, nil
}
