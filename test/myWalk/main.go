package main

import (
	"log"
	"os"
)

func myWalk(currentPath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		dir, err := os.Open(currentPath)
		if err != nil {
			return err
		}
		defer dir.Close()

		fis, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, fi := range fis {
			if fi.Name() != "." && fi.Name() != ".." {
				myWalk(currentPath+"/"+fi.Name(), fi, err)
			}
		}
	} else {
		log.Printf("myWalk %s [%d] %t", currentPath, info.Size(), info.Mode()&os.ModeSymlink != 0)
	}

	return nil
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		info, err := os.Lstat(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}
		if err := myWalk(os.Args[i], info, nil); err != nil {
			log.Fatal(err)
		}
	}
}
