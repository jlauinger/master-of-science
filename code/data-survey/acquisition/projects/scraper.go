package projects

import (
	"context"
	"data-acquisition/lexical"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/gocarina/gocsv"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetProjects(dataDir, downloadDir string, download, createForks bool, forkTargetOrg, accessToken string) {
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Printf("Saving project data to %s\n", projectsFilename)

	headerWritten := false
	projectsFile, err := os.OpenFile(projectsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer projectsFile.Close()

	fmt.Println("Getting information about top 500 Go projects...")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

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
			path := downloadDir + "/" + *repo.Name
			revision := ""

			fmt.Printf("%v. %v\n", (page-1)*100+(i+1), *repo.CloneURL)

			if download {
				revision = downloadRepo(repo, path)
			}

			if createForks {
				createFork(client, repo, forkTargetOrg)
			}

			project := lexical.ProjectData{
				Rank:           i + 1,
				Name:           repo.GetFullName(),
				GithubCloneUrl: repo.GetCloneURL(),
				NumberOfStars:  repo.GetStargazersCount(),
				NumberOfForks:  repo.GetForksCount(),
				GithubId:       *repo.ID,
				Revision:       revision,
				CreatedAt:      lexical.DateTime{Time: repo.CreatedAt.Time},
				LastPushedAt:   lexical.DateTime{Time: repo.PushedAt.Time},
				UpdatedAt:      lexical.DateTime{Time: repo.UpdatedAt.Time},
				Size:           *repo.Size,
				CheckoutPath:   path,
			}

			if headerWritten {
				_ = gocsv.MarshalWithoutHeaders([]lexical.ProjectData{project}, projectsFile)
			} else {
				headerWritten = true
				_ = gocsv.Marshal([]lexical.ProjectData{project}, projectsFile)
			}
		}
	}
}

func createFork(client *github.Client, repo github.Repository, targetOrg string) {
	_, _, err := client.Repositories.CreateFork(context.Background(), *repo.Owner.Name, *repo.Name, &github.RepositoryCreateForkOptions{
		Organization: targetOrg,
	})
	_, ok := err.(*github.AcceptedError)
	if !ok && err != nil {
		fmt.Printf("ERROR: %v!", err)
	}
	fmt.Printf("  forked to %s/%s\n", targetOrg, *repo.Name)
}

func downloadRepo(repo github.Repository, path string) string {
	fmt.Printf("  Downloading to %v ...", path)

	cloneCtx, err := git.PlainClone(path, false, &git.CloneOptions{
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

	head, err := cloneCtx.Head()
	if err != nil {
		fmt.Printf("ERROR: %v!", err)
	}
	revision := head.Hash().String()

	fmt.Printf("  checked out revision %s\n", revision)
	fmt.Println("  done!")

	return revision
}
