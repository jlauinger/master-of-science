package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	downloadPath := "./download"

	fmt.Println("Getting top 500 Go projects...")

	client := github.NewClient(nil)

	for page := 1; page <= 5; page++ {
		repos, _, err := client.Search.Repositories(context.Background(), "language:Go", &github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page: page,
			},
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		for i, repo := range repos.Repositories {
			path := downloadPath + "/" + *repo.Name

			fmt.Printf("%v. %v\n", (page-1)*100+(i+1), *repo.CloneURL)
			fmt.Printf("  Downloading to %v ...", path)

			_, err := git.PlainClone(path, false, &git.CloneOptions{
				URL:               *repo.CloneURL,
				Depth:             1,
				Progress:          nil,
			})

			if err != nil {
				fmt.Printf("ERROR: %v!\n", err)
			} else {
				fmt.Printf("done\n")
			}
		}
	}
}
