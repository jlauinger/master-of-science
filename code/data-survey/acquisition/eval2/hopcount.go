package eval2

import "github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"

func analyzeDepTree(packages []*lexical.PackageData) {
	packagesGetImported := make(map[string]bool, len(packages))
	packagesMap := make(map[string]*lexical.PackageData, len(packages))

	for _, pkg := range packages {
		packagesGetImported[pkg.ImportPath] = false
		packagesMap[pkg.ImportPath] = pkg
	}

	for _, pkg := range packages {
		for _, childPath := range pkg.Imports {
			if childPath == "C" {
				continue
			}
			child := packagesMap[childPath]
			packagesGetImported[child.ImportPath] = true
		}
	}

	rootPackages := make([]*lexical.PackageData, 0)

	for pkgPath, getsImported := range packagesGetImported {
		if getsImported {
			continue
		}
		pkg := packagesMap[pkgPath]
		if pkg.ImportPath == "runtime/cgo" {
			continue
		}

		rootPackages = append(rootPackages, pkg)
	}

	analyzeHopCountBFS(rootPackages, packagesMap)
}

func analyzeHopCountBFS(rootPackages []*lexical.PackageData, packagesMap map[string]*lexical.PackageData) {

	type PackageAndPotentialHopCount struct {
		PotentialHopCount int
		ImportStack       []string
		Pkg               *lexical.PackageData
	}

	queue := make([]PackageAndPotentialHopCount, 0)
	seen := make(map[string]bool, 0)

	for _, rootPkg := range rootPackages {
		queue = append(queue, PackageAndPotentialHopCount{
			PotentialHopCount: 0,
			Pkg:               rootPkg,
		})
	}

	var queueItem PackageAndPotentialHopCount
	for {
		if len(queue) == 0 {
			// finished
			break
		}

		// shift first item
		queueItem, queue = queue[0], queue[1:]

		// set hop count and mark package as seen in BFS
		queueItem.Pkg.HopCount = queueItem.PotentialHopCount
		queueItem.Pkg.ImportStack = queueItem.ImportStack
		seen[queueItem.Pkg.ImportPath] = true

		for _, childPath := range queueItem.Pkg.Imports {
			// ignore the Cgo import
			if childPath == "C" {
				continue
			}
			child, ok := packagesMap[childPath]
			if !ok {
				panic("child not found")
			}
			_, ok = seen[child.ImportPath]
			if !ok {
				// push unseen children to the back of the queue
				queue = append(queue, PackageAndPotentialHopCount{
					PotentialHopCount: queueItem.PotentialHopCount + 1,
					ImportStack:       append(queueItem.ImportStack, child.ImportPath),
					Pkg:               child,
				})
			}
		}
	}
}
