# gfdupes cli

```
g = golang
f = file
dupes = duplicates
```

gfdupes is a cli written in golang used to find duplicate files. It draws inspiration from fdupes.

### Dev Notes

Uses golang concurrency pipeline pattern.

```
generate file list ->
gather sizes and modes -> 
hash duplicates -> 
distill duplicates -> 
print
```

### TODO

* progress bar
* buffer size
* cobra cli
* check if paths overlap
* follow symlinks
* consider hardlink in diff
* delete 2-n diffs
* gunzip content
* ntfs alternate file streams
* -exec option
* -verbose
* -version

### References

* https://golang.org
* https://blog.golang.org/pipelines
* https://github.com/adrianlopezroche/fdupes
