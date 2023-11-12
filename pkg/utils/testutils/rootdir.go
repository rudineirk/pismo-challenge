package testutils

import (
	"os"
	"strings"
)

func GetRootDir() (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dirs := strings.Split(dir, "/")
	rootDir := ""

	for _, dir := range dirs {
		rootDir += dir + "/"

		info, _ := os.Stat(rootDir + "go.mod")

		if info != nil {
			if info.Name() == "go.mod" {
				break
			}
		}
	}

	return rootDir, nil
}

func SetRootCwd() error {
	rootDir, err := GetRootDir()
	if err != nil {
		return err
	}

	return os.Chdir(rootDir)
}
