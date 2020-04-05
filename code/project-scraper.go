package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
				PerPage: 100, // 100,
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
				fmt.Println("done")
			}

			fmt.Printf("  Vendoring Go modules ...")

			var goModPaths []string

			err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err == nil && strings.ToLower(info.Name()) == "go.mod" {
					goModPaths = append(goModPaths, path[:len(path)-len("go.mod")])
				}
				return nil
			})
			if err != nil {
				fmt.Printf("ERROR: %v!\n", err)
			}

			for _, goModPath := range goModPaths {
				fmt.Printf("\n  Running go mod vendor in %v ...", goModPath)

				cmd := exec.Command("go", "mod", "vendor")
				cmd.Dir = goModPath

				err = cmd.Run()
				if err != nil {
					fmt.Printf("ERROR: %v!", err)
				} else {
					fmt.Printf("done")
				}
			}

			fmt.Println("  done!")
		}
	}
}
