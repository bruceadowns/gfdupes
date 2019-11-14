package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func getShellScript(rootpath string) []string {
	list := make([]string, 0)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sh" {
			list = append(list, path)
		}
		return nil
	})
	if err != nil {
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
