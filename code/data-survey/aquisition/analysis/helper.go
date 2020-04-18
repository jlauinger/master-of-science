package analysis

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func getRegistryFromImportPath(importPath string) string {
	pathComponents := strings.Split(importPath, "/")

	if len(pathComponents) <= 1 {
		return "std"
	}

	var registryComponents []string
	if pathComponents[1] == "x" {
		registryComponents = pathComponents[0:2]
	} else {
		registryComponents = pathComponents[0:1]
	}

	return strings.Join(registryComponents, "/")
}

func copyFiles(files map[string]string, copyDestination string) {
	fmt.Println("  copying potentially vulnerable files...")

	for srcFilename, destFilename := range files {
		fullDestFilename := fmt.Sprintf("%s/%s", copyDestination, destFilename)
		
		src, err := os.Open(srcFilename)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "copy-src",
				ProjectName:       "",
				PackageImportPath: "",
				FileName:          srcFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		destFolder := path.Dir(fullDestFilename)
		err = os.MkdirAll(destFolder, 0755)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "copy-mkdir",
				ProjectName:       "",
				PackageImportPath: fullDestFilename,
				FileName:          srcFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		dest, err := os.OpenFile(fullDestFilename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "copy-dest",
				ProjectName:       "",
				PackageImportPath: fullDestFilename,
				FileName:          srcFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		_, err = io.Copy(dest, src)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "copy-copy",
				ProjectName:       "",
				PackageImportPath: fullDestFilename,
				FileName:          srcFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		src.Close()
		dest.Close()
	}
}