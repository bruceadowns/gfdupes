package main

import (
	"log"
	"os"
	"path/filepath"
)

func addWildcard(s string) string {
	if len(s) == 0 {
		return "*"
	}
	if s[len(s)-1] == '/' {
		return s + "*"
	}
	return s + "/*"
}

func myGlob(currentPath string) error {
	files, err := filepath.Glob(addWildcard(currentPath))
	if err != nil {
		return err
	}

	for _, v := range files {
		info, err := os.Lstat(v)
		if err != nil {
			log.Fatal(err)
		}

		if !info.IsDir() {
			log.Printf("myGlob: %s", v)
		}
	}

	return nil
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		if err := myGlob(os.Args[i]); err != nil {
			log.Fatal(err)
		}
	}
}
