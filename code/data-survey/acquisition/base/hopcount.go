package base

import (
	"fmt"
)

/**
 * identifies the minimal depth in the dependency tree for each package and updates it into the hop count field
 */
func analyzeHopCounts(packages []*PackageData) []*PackageData {
	fmt.Println("  analyzing dependency structure and hop counts...")

	// initialize a hash set of packages by their import path flagging whether they get imported somewhere, and a hash
	// map to map from the import path back to the package
	packagesGetImported := make(map[string]bool, len(packages))
	packagesMap := make(map[string]*PackageData, len(packages))

	// go through all packages and fill those data structures
	for _, pkg := range packages {
		packagesGetImported[pkg.ImportPath] = false
		packagesMap[pkg.ImportPath] = pkg
	}

	// then to through all packages again
	for _, pkg := range packages {
		// go through all the imports (this does not have to be done recursively because in the outer loop we are already
		// covering all packages anyway)
		for _, childPath := range pkg.Imports {
			// skip the C pseudo package
			if childPath == "C" {
				continue
			}
			// identify the package and set it to seen
			child := packagesMap[childPath]
			packagesGetImported[child.ImportPath] = true
		}
	}

	// initialize a list for all root packages, i.e. those that do not get imported by other packages
	rootPackages := make([]*PackageData, 0)

	// to fill it, go through the hash set of packages populated before
	for pkgPath, getsImported := range packagesGetImported {
		// if the package got imported, it is not a root package and should be skipped
		if getsImported {
			continue
		}
		// otherwise, identify the package using the hash map
		pkg := packagesMap[pkgPath]
		// if it is the special runtime/cgo package, skip it
		if pkg.ImportPath == "runtime/cgo" {
			continue
		}
		// finally, add it to the list of root packages
		rootPackages = append(rootPackages, pkg)
	}

	// use breadth first search to identify the import depths
	analyzeHopCountBFS(rootPackages, packagesMap)

	return rootPackages
}

/**
 * identifies the import depth / hop count for all packages reachable through the root packages using breadth first
 * search. Using BFS is important because there are many duplicate paths and using DFS will be ridiculously slow
 */
func analyzeHopCountBFS(rootPackages []*PackageData, packagesMap map[string]*PackageData) {

	// this is a structure holding the import stack and hop count as well as the package itself and can be inserted
	// into the BFS working queue to remember the hop count this package would get if it is drawn from the queue
	type PackageAndPotentialHopCount struct {
		PotentialHopCount int
		ImportStack       []string
		Pkg               *PackageData
	}

	// initialize a BFS working queue and a hash set of packages that have already been seen. Since we are interested
	// in the minimal import depth, BFS is a perfect match here because once a package has been seen once, it can not
	// be on a lower depth and therefore does not have to be analyzed ever again
	queue := make([]PackageAndPotentialHopCount, 0)
	seen := make(map[*PackageData]bool, 0)

	// append all root packages at depth 0 to the BFS queue
	for _, rootPkg := range rootPackages {
		queue = append(queue, PackageAndPotentialHopCount{
			PotentialHopCount: 0,
			Pkg:               rootPkg,
		})
	}

	// draw from the queue until BFS terminates
	var queueItem PackageAndPotentialHopCount
	for {
		// terminate if the queue is empty
		if len(queue) == 0 {
			// finished
			break
		}

		// shift first item
		queueItem, queue = queue[0], queue[1:]

		// set hop count and mark package as seen in BFS. It is important to do this before analyzing the children
		queueItem.Pkg.HopCount = queueItem.PotentialHopCount
		queueItem.Pkg.ImportStack = queueItem.ImportStack
		seen[queueItem.Pkg] = true

		// then go through all of the package children, i.e. the packages imported by this package
		for _, childPath := range queueItem.Pkg.Imports {
			// ignore the Cgo import
			if childPath == "C" {
				continue
			}
			// identify the child package
			child, ok := packagesMap[childPath]
			if !ok {
				panic("child not found")
			}
			// check if the child has previously been seen
			_, ok = seen[child]
			if !ok {
				// push unseen child to the back of the queue and mark it seen
				queue = append(queue, PackageAndPotentialHopCount{
					PotentialHopCount: queueItem.PotentialHopCount + 1,
					ImportStack:       append(queueItem.ImportStack, child.ImportPath),
					Pkg:               child,
				})
				seen[child] = true
			}
		}
	}
}
