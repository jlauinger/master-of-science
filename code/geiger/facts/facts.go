package facts

import (
	"go/types"
	"sync"
)

type PackageInfo struct {
	ThisCount   int
}

var mu *sync.Mutex
var infos map[*types.Package]*PackageInfo

func Init() {
	infos = make(map[*types.Package]*PackageInfo, 0)
	mu = new(sync.Mutex)
}

func Store(pkg *types.Package, thisCount int) {
	mu.Lock()
	infos[pkg] = &PackageInfo{
		ThisCount:  thisCount,
	}
	mu.Unlock()
}

func GetAll() map[*types.Package]*PackageInfo {
	return infos
}
