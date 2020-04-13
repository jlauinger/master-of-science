package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/go-git/go-git/v5"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	downloadPath := "/home/johannes/studium/s14/masterarbeit/download"

	fmt.Println("Getting top 500 Go projects...")

	client := github.NewClient(nil)

	file, err := os.Create("projects.csv")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"rank", "name", "github_clone_url", "number_of_stars", "number_of_forks",
		"github_id", "created_at", "last_pushed_at", "updated_at", "size", "checkout_path"})

	for page := 1; page <= 1; page++ {
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

			writer.Write([]string{
				strconv.Itoa(i + 1),
				repo.GetFullName(),
				repo.GetCloneURL(),
				strconv.Itoa(repo.GetStargazersCount()),
				strconv.Itoa(repo.GetForksCount()),
				strconv.FormatInt(*repo.ID, 10),
				repo.CreatedAt.String(),
				repo.PushedAt.String(),
				repo.UpdatedAt.String(),
				strconv.Itoa(*repo.Size),
				path,
			})

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
