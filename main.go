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

var theVersion = "v0.1"

type args struct {
	bufferSize  int
	hashType    string
	noEmpty     bool
	noHidden    bool
	perm        bool
	recurse     bool
	showVersion bool
	verbose     bool
}

func newArgs() (res args) {
	flag.IntVar(&res.bufferSize, "buffer", 0, "buffer size used for channel pipeline")
	flag.StringVar(&res.hashType, "hash", "md5", "hash type of md5 or sha1")
	flag.BoolVar(&res.noEmpty, "noempty", false, "exclude empty files in difference")
	flag.BoolVar(&res.noHidden, "nohidden", false, "exclude hidden files in difference")
	flag.BoolVar(&res.perm, "perm", false, "consider permissions in difference")
	flag.BoolVar(&res.recurse, "recurse", false, "recurse")
	flag.BoolVar(&res.showVersion, "version", false, "show version")
	flag.BoolVar(&res.verbose, "verbose", false, "debug logging to stderr")
	flag.Parse()

	switch res.hashType {
	case "md5", "sha", "sha1":
	default:
		log.Fatalf("invalid hash type: %s expect md5 or sha1", res.hashType)
	}

	return
}

var fnLog func(format string, v ...interface{}) = nil

func vLog(format string, v ...interface{}) {
	if fnLog != nil {
		fnLog(format, v...)
	}
}

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

func computeHash(filename, hashType string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		vLog("Error opening: %s [%s]", filename, err)
		return "", err
	}
	defer f.Close()

	var h hash.Hash
	switch hashType {
	case "md5":
		h = md5.New()
	case "sha", "sha1":
		h = sha1.New()
	default:
		return "", fmt.Errorf("invalid hash type: %s", hashType)
	}

	if _, err = io.Copy(h, f); err != nil {
		vLog("Error hashing file: %s [%s]", filename, err)
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum), nil
}

func genFilesRecursive(paths []string, bufferSize int, noEmpty, noHidden bool) (res chan fileWalk) {
	res = make(chan fileWalk, bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		currentPath := paths[i]
		vLog("Walk path: %s", currentPath)

		go func() {
			defer wg.Done()

			if err := filepath.Walk(currentPath, func(path string, info os.FileInfo, errIn error) error {
				if errIn != nil {
					vLog("Error walking path: %s", errIn)
					return errIn
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

func addWildCard(s string) string {
	if len(s) == 0 {
		return "*"
	}

	if s[len(s)-1] == '/' {
		return s + "*"
	}

	return s + "/*"
}

func genFiles(paths []string, bufferSize int, noEmpty, noHidden bool) (res chan fileWalk) {
	res = make(chan fileWalk, bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		wg.Add(1)
		path := addWildCard(paths[i])

		go func() {
			defer wg.Done()

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

func gatherFiles(ch <-chan fileWalk, bufferSize int, perm bool) (res chan fileAttrs) {
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

			files := fam[fa]
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

func hashFiles(ch <-chan fileAttrs, bufferSize int, hashType string) (res chan fileAttrsExt) {
	res = make(chan fileAttrsExt, bufferSize)

	var wg sync.WaitGroup
	for f := range ch {
		wg.Add(1)
		f1 := f

		go func() {
			defer wg.Done()

			if h, err := computeHash(f1.name, hashType); err == nil {
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

func distillFiles(ch <-chan fileAttrsExt, bufferSize int) (res chan []string) {
	res = make(chan []string, bufferSize)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		faem := make(fileAttrsExtMap)
		for f := range ch {
			fae := fileAttrsExt{size: f.size, mode: f.mode, hash: f.hash}
			files := faem[fae]
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
	for f := range ch {
		if once {
			fmt.Println()
		}
		once = true

		for _, v := range f {
			fmt.Println(v)
		}
	}
}

func main() {
	args := newArgs()
	if args.showVersion {
		fmt.Println(theVersion)
	}
	if args.verbose {
		fnLog = log.Printf
	}

	// pipeline

	var genChan chan fileWalk
	if args.recurse {
		genChan = genFilesRecursive(flag.Args(), args.bufferSize, args.noEmpty, args.noHidden)
	} else {
		genChan = genFiles(flag.Args(), args.bufferSize, args.noEmpty, args.noHidden)
	}

	gatherChan := gatherFiles(genChan, args.bufferSize, args.perm)

	hashChan := hashFiles(gatherChan, args.bufferSize, args.hashType)

	distillChan := distillFiles(hashChan, args.bufferSize)

	printFilenames(distillChan)
}
