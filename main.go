package main

import (
	"fmt"
	"nodenv-lts-install/src/lts_install"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	cwd := os.Args[1]
	versionQuery := os.Args[2]
	dryRun := len(os.Args) > 3

	if err := os.MkdirAll(cwd, os.ModePerm); err != nil {
		panic(err)
		return
	}

	if _, err := os.Stat(filepath.Join(cwd, versionQuery)); os.IsExist(err) {
		return
	}

	version, err := lts_install.GetVersion(versionQuery)

	if err != nil {
		panic(err)
		return
	}

	if version == nil {
		return
	}

	hash, err := lts_install.GetHash(*version, osName, arch)

	if err != nil {
		panic(err)
		return
	}

	fmt.Println(version.Version[1:])

	if _, err = os.Stat(filepath.Join(cwd, version.Version[1:])); os.IsExist(err) {
		return
	}

	fileContent := lts_install.MakeVersionFileContent(*version, hash, osName, arch)

	if dryRun {
		fmt.Println(string(fileContent))
		return
	}

	err = os.WriteFile(
		filepath.Join(cwd, version.Version[1:]),
		fileContent,
		0644,
	)

	if err != nil {
		panic(err)
	}
}
