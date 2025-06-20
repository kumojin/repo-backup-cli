package main

import (
	"context"

	"github.com/google/go-github/v72/github"
)

func main() {
	client := github.NewClient(nil).WithAuthToken("your_token_here")

	repos, _, err := client.Repositories.ListByOrg(context.TODO(), "github", &github.RepositoryListByOrgOptions{
		Type: "private",
	})
	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		// Do something with each repository
		println(*repo.Name)
	}
}
