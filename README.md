# gfdupes cli

```
g = go
f = file
dupes = duplicates
```

gfdupes is a cli written in golang used to find duplicate files. It draws inspiration from fdupes.

## Dev Notes

Uses golang concurrency pipeline pattern.

```
generate file list ->
gather sizes and modes -> 
hash duplicates -> 
distill duplicates -> 
print
```

## References

* https://golang.org
* https://github.com/adrianlopezroche/fdupes
