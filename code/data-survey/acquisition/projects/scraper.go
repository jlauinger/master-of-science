package projects

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/**
 * identifies the Top 500 most starred open-source Go projects, forks and downloads them
 */
func GetProjects(datasetSize int, dataDir, downloadDir string, download, createForks bool, accessToken string) {
	// build the projects CSV filename from the configuration
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Printf("Saving project data to %s\n", projectsFilename)

	if err := base.OpenProjectsFile(projectsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	fmt.Printf("Getting information about top %d Go projects...\n", datasetSize)

	// set up the Github SDK client with the access token as provided by the configuration. Querying projects can be
	// done without authentication, but forking projects does require an access token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	// split up the projects into pages of 100 each because Github's rate limiting does not allow more
	pages := datasetSize / 100
	for page := 1; page <= pages + 1; page++ {
		// calculate the size of this page, if needed at all
		pageSize := base.Min(datasetSize - ((page-1)*100), 100)
		if pageSize <= 0 {
			continue
		}

		// search for repositories with language Go on Github. They are ordered by stars automatically
		repos, _, err := client.Search.Repositories(context.Background(), "language:Go", &github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: pageSize,
				Page: page,
			},
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// got through all the repositories in the search results
		for i, repo := range repos.Repositories {
			// build the download path
			path := downloadDir + "/" + repo.GetFullName()
			revision := ""

			fmt.Printf("%v. %v\n", (page-1)*100+(i+1), *repo.CloneURL)

			// download and checkout the project if requested
			if download {
				revision = downloadRepo(repo, path)
			}

			// create a fork if requested
			if createForks {
				createFork(client, repo)
			}

			// set up the project data to be written into CSV
			project := base.ProjectData{
				Rank:           (page-1)*100+(i+1),
				Name:           repo.GetFullName(),
				GithubCloneUrl: repo.GetCloneURL(),
				NumberOfStars:  repo.GetStargazersCount(),
				NumberOfForks:  repo.GetForksCount(),
				GithubId:       *repo.ID,
				Revision:       revision,
				CreatedAt:      base.DateTime{Time: repo.CreatedAt.Time},
				LastPushedAt:   base.DateTime{Time: repo.PushedAt.Time},
				UpdatedAt:      base.DateTime{Time: repo.UpdatedAt.Time},
				Size:           *repo.Size,
				CheckoutPath:   path,
			}

			// write the project into the CSV file
			err = base.WriteProject(project)
			if err != nil {
				panic(err)
			}
		}
	}
}

/**
 * forks the given repository into the account identified by the Github access token given by configuration
 */
func createFork(client *github.Client, repo github.Repository) {
	// extract the repository owner account as the part before the slash
	components := strings.Split(repo.GetFullName(), "/")
	owner := components[0]

	// create a fork of the repository using the Github API
	_, _, err := client.Repositories.CreateFork(context.Background(), owner, *repo.Name,
		&github.RepositoryCreateForkOptions{})
	_, ok := err.(*github.AcceptedError)
	if !ok && err != nil {
		fmt.Printf("ERROR: %v!", err)
	}

	fmt.Printf("  forked to %s\n", *repo.Name)
}

/**
 * downloads, checks out the given repository and returns the current revision SHA. Additionally, all go.mod files
 * are vendored to ensure the dependency modules are downloaded
 */
func downloadRepo(repo github.Repository, path string) string {
	fmt.Printf("  Downloading to %v ...", path)

	// clone the repository into the download path using the Git library. Since I am not interested in the history,
	// a shallow clone is enough
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

	// this is a list of go.mod files contained in this repository
	var goModPaths []string

	// walk through the directories contained in the repository
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// is the current file a go.mod file? There can be several per repository
		if err == nil && strings.ToLower(info.Name()) == "go.mod" {
			// if so, append the directory of this go.mod file to the paths list
			goModPaths = append(goModPaths, path[:len(path)-len("go.mod")])
		}
		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: %v!\n", err)
	}

	// go through all the go.mod paths
	for _, goModPath := range goModPaths {
		fmt.Printf("\n  Running go mod vendor in %v ...", goModPath)

		// build the go mod vendor command for this go.mod file
		cmd := exec.Command("go", "mod", "vendor")
		cmd.Dir = goModPath

		// and run it to ensure dependencies are properly downloaded
		err = cmd.Run()
		if err != nil {
			fmt.Printf("ERROR: %v!", err)
		} else {
			fmt.Printf("done")
		}
	}

	// identify the revision of this checkout using the Git library
	head, err := cloneCtx.Head()
	if err != nil {
		fmt.Printf("ERROR: %v!", err)
	}
	revision := head.Hash().String()

	fmt.Printf("\n  done with revision %s\n", revision)

	return revision
}
