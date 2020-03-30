package main

import (
	"fmt"
	"github.com/KyleBanks/depth"
)

func main() {
	fmt.Println("Show dependencies of a Go project")

	var t depth.Tree

	err := t.Resolve("github.com/KyleBanks/depth")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("'%v' has %v dependencies:\n", t.Root.Name, len(t.Root.Deps))

	for _, dep := range t.Root.Deps {
		fmt.Printf("  %v\n", dep.Name)
	}
}
