package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func getShellScript(rootpath string) []string {
	var list []string

	if err := filepath.Walk(rootpath, func(path string, info os.FileInfo, errIn error) error {
		if errIn != nil {
			return errIn
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".sh" {
			list = append(list, path)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return list
}

func main() {
	gopath := os.Getenv("GOPATH")
	log.Printf("[%s/bin]", gopath)

	list := getShellScript(gopath)
	for i, p := range list {
		log.Printf("%d: %s/%s", i, path.Dir(p), path.Base(p))
	}
}
