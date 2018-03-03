package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type fileWalk struct {
	name string
	fi   os.FileInfo
}

type fileAttrs struct {
	name string
	size int64
	mode os.FileMode
}
type fileAttrsMap map[fileAttrs][]string

type fileAttrsExt struct {
	name string
	size int64
	mode os.FileMode
	hash string
}
type fileAttrsExtMap map[fileAttrsExt][]string

func computeHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening: %s [%s]", filename, err)
		return "", err
	}
	defer file.Close()

	//h := sha1.New()
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Printf("Error hashing file: %s [%s]", filename, err)
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum), nil
}

func genFiles(paths []string) (res chan fileWalk) {
	res = make(chan fileWalk)

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		currentPath := paths[i]

		go func() {
			defer wg.Done()

			if err := filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Printf("Error walking path: %s", err)
					return nil
				}

				if info.IsDir() {
					return nil
				}

				res <- fileWalk{name: path, fi: info}
				return nil
			}); err != nil {
				log.Printf("Error walking file names: %s", err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return
}

func gatherFiles(ch <-chan fileWalk) (res chan fileAttrs) {
	res = make(chan fileAttrs)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		fam := make(fileAttrsMap)
		for f := range ch {
			fa := fileAttrs{size: f.fi.Size(), mode: f.fi.Mode()}
			files, ok := fam[fa]
			if !ok {
				files = make([]string, 0)
			}
			fam[fa] = append(files, f.name)
		}

		for k, v := range fam {
			if len(v) > 1 {
				for _, vv := range v {
					res <- fileAttrs{name: vv, size: k.size, mode: k.mode}
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(res)
	}()

	return
}

func hashFiles(ch <-chan fileAttrs) (res chan fileAttrsExt) {
	res = make(chan fileAttrsExt)

	var wg sync.WaitGroup
	for f := range ch {
		wg.Add(1)
		f1 := f

		go func() {
			defer wg.Done()

			h, err := computeHash(f1.name)
			if err == nil {
				res <- fileAttrsExt{name: f1.name, size: f1.size, mode: f1.mode, hash: h}
			} else {
				log.Printf("Error hashing %s: %s", f1.name, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return
}

func distillFiles(ch <-chan fileAttrsExt) (res chan []string) {
	res = make(chan []string)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		faem := make(fileAttrsExtMap)
		for f := range ch {
			fae := fileAttrsExt{size: f.size, mode: f.mode, hash: f.hash}
			files, ok := faem[fae]
			if !ok {
				files = make([]string, 0)
			}
			faem[fae] = append(files, f.name)
		}

		for _, v := range faem {
			if len(v) > 1 {
				res <- v
			}
		}
	}()

	go func() {
		wg.Wait()
		close(res)
	}()

	return
}

func printFilenames(ch <-chan []string) {
	once := false
	for filesets := range ch {
		if once {
			fmt.Println()
		}
		once = true

		for _, v := range filesets {
			fmt.Println(v)
		}
	}
}

func main() {
	// main pipeline
	genChan := genFiles(os.Args[1:])
	gatherChan := gatherFiles(genChan)
	hashChan := hashFiles(gatherChan)
	distillChan := distillFiles(hashChan)
	printFilenames(distillChan)
}

/*
	command line (cobra)

	arguments
	* gfdupes [dir1] [dir2]
	* check if paths overlap

	flags
	* version
	* verbose
	* progress bar
	* recurse
	* follow symlinks
	* consider permissions in diff
	* consider hardlink in diff
	* delete 2-n diffs
	* gunzip content
	* ntfs alternate file streams
	* hash type md5/sha/etc
	* provide -exec option for

	recursive behavior
	* filepath.Glob for non-recursive behavior
	* filepath.Walk for recursive behavior
*/
