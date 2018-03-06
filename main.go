package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
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

func vLog(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v)
	}
}

func computeHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		vLog("Error opening: %s [%s]", filename, err)
		return "", err
	}
	defer file.Close()

	var h hash.Hash
	switch hashType {
	case "md5":
		h = md5.New()
	case "sha", "sha1":
		h = sha1.New()
	default:
		log.Fatalf("Invalid hash type: %s", hashType)
	}
	if _, err := io.Copy(h, file); err != nil {
		vLog("Error hashing file: %s [%s]", filename, err)
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum), nil
}

func genFilesRecursive(paths []string) (res chan fileWalk) {
	res = make(chan fileWalk, bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		currentPath := paths[i]
		vLog("Walk path: %s", currentPath)

		go func() {
			defer wg.Done()

			if err := filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					vLog("Error walking path: %s", err)
					return nil
				}

				vLog("Found file: %s", path)

				if info.IsDir() {
					vLog("Ignore directory: %s", path)
					return nil
				}

				if noEmpty && info.Size() == 0 {
					vLog("Ignore empty file: %s", path)
					return nil
				}

				if noHidden && info.Name()[0] == '.' {
					vLog("Ignore hidden file: %s", path)
					return nil
				}

				res <- fileWalk{name: path, fi: info}
				return nil
			}); err != nil {
				vLog("Error walking file names: %s", err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return
}

func genFiles(paths []string) (res chan fileWalk) {
	res = make(chan fileWalk, bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		currentPath := paths[i]

		go func() {
			defer wg.Done()

			path := func(s string) string {
				if len(s) == 0 {
					return "*"
				}
				if s[len(s)-1] == '/' {
					return s + "*"
				}
				return s + "/*"
			}(currentPath)

			vLog("Globbing path: %s", path)
			files, err := filepath.Glob(path)
			if err != nil {
				vLog("Error globbing %s: %s", path, err)
				return
			}

			for _, v := range files {
				vLog("Found file: %s", v)

				info, err := os.Lstat(v)
				if err != nil {
					log.Fatal(err)
				}

				if info.IsDir() {
					vLog("Ignore directory: %s", v)
					continue
				}

				if noEmpty && info.Size() == 0 {
					vLog("Ignore empty file: %s", v)
					continue
				}

				if noHidden && info.Name()[0] == '.' {
					vLog("Ignore hidden file: %s", v)
					continue
				}

				res <- fileWalk{name: v, fi: info}
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
	res = make(chan fileAttrs, bufferSize)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		fam := make(fileAttrsMap)
		for f := range ch {
			fa := fileAttrs{size: f.fi.Size()}
			if perm {
				fa.mode = f.fi.Mode()
			} else {
				vLog("Do not consider permissions in difference for: %s", f.name)
			}

			files, ok := fam[fa]
			if !ok {
				files = make([]string, 0)
			}
			fam[fa] = append(files, f.name)
		}

		for k, v := range fam {
			if len(v) > 1 {
				vLog("Found initial duplicate file: %s", v)

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
	res = make(chan fileAttrsExt, bufferSize)

	var wg sync.WaitGroup
	for f := range ch {
		wg.Add(1)
		f1 := f

		go func() {
			defer wg.Done()

			h, err := computeHash(f1.name)
			if err == nil {
				vLog("%s hash for %s: %s", hashType, f1.name, h)
				res <- fileAttrsExt{name: f1.name, size: f1.size, mode: f1.mode, hash: h}
			} else {
				vLog("Error %s hashing %s: %s", hashType, f1.name, err)
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
	res = make(chan []string, bufferSize)

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
				vLog("Found hashed duplicate file: %s", v)
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

var bufferSize int
var hashType string
var noEmpty bool
var noHidden bool
var perm bool
var recurse bool
var showVersion bool
var verbose bool

func init() {
	flag.IntVar(&bufferSize, "buffer", 0, "buffer size used for channel pipeline")
	flag.StringVar(&hashType, "hash", "md5", "hash type of md5 or sha1")
	flag.BoolVar(&noEmpty, "noempty", false, "exclude empty files in difference")
	flag.BoolVar(&noHidden, "nohidden", false, "exclude hidden files in difference")
	flag.BoolVar(&perm, "perm", false, "consider permissions in difference")
	flag.BoolVar(&recurse, "recurse", false, "recurse")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&verbose, "verbose", false, "debug logging to stderr")
	flag.Parse()

	switch hashType {
	case "md5", "sha", "sha1":
	default:
		log.Fatalf("Invalid hash type: %s. Expect md5 or sha1.", hashType)
	}
}

var version = "v0.1"

func main() {
	if showVersion {
		fmt.Println(version)
	}

	// main pipeline

	var genChan chan fileWalk
	if recurse {
		genChan = genFilesRecursive(flag.Args())
	} else {
		genChan = genFiles(flag.Args())
	}

	gatherChan := gatherFiles(genChan)

	hashChan := hashFiles(gatherChan)

	distillChan := distillFiles(hashChan)

	printFilenames(distillChan)
}
