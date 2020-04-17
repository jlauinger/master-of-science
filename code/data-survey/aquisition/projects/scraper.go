package projects

import (
	"context"
	"data-aquisition/analysis"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/gocarina/gocsv"
	"github.com/google/go-github/github"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetProjects(dataDir, downloadDir string, download bool) {
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
			path := downloadDir + "/" + *repo.Name

			fmt.Printf("%v. %v\n", (page-1)*100+(i+1), *repo.CloneURL)
			
			project := analysis.ProjectData{
				ProjectRank:           i + 1,
				ProjectName:           repo.GetFullName(),
				ProjectGithubCloneUrl: repo.GetCloneURL(),
				ProjectNumberOfStars:  repo.GetStargazersCount(),
				ProjectNumberOfForks:  repo.GetForksCount(),
				ProjectGithubId:       *repo.ID,
				ProjectCreatedAt:      analysis.DateTime{Time: repo.CreatedAt.Time},
				ProjectLastPushedAt:   analysis.DateTime{Time: repo.PushedAt.Time},
				ProjectUpdatedAt:      analysis.DateTime{Time: repo.UpdatedAt.Time},
				ProjectSize:           *repo.Size,
				ProjectCheckoutPath:   path,
			}

			if headerWritten {
				gocsv.MarshalWithoutHeaders([]analysis.ProjectData{project}, projectsFile)
			} else {
				gocsv.Marshal([]analysis.ProjectData{project}, projectsFile)
			}

			if download {
				downloadRepo(repo, path)
			}
		}
	}
}

func downloadRepo(repo github.Repository, path string) {
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
