# gfdupes cli

```
g       golang port of
f       file
dupes   duplicates (fdupes)
c       command
l       line
i       interface
```

gfdupes is a cli written in golang used to find duplicate files. It draws inspiration from fdupes.

### CLI

##### Implemented

* ```-buffer``` - buffer size for pipeline channel
* ```-hash``` - use {md5,sha1} hash
* ```-perm``` - consider permissions in difference
* ```-noempty``` - exclude empty files
* ```-nohidden``` - exclude hidden files
* ```-recurse``` - recurse subdirectories
* ```-verbose``` - debug logging to stderr
* ```-version``` - show version

##### Ideas

* ```-size +10m -size -20m``` - multiple file size constraints [kmgtpKMGTP]
* ```[!] [-not] -name [glob] -name ...``` - multiple basename wildcards
* ```-maxdepth``` - recurse to a maximum depth
* ```-mtime``` - modified in last n [smhdwy] units
* ```-xdev``` - prevent recursion across devices
* ```-symlinks``` - follow symlinks
* ```-hardlinks``` - hardlink in diff
* ```-delete``` - prompt to delete
* ```-gzip``` - consider gzip uncompressed file stats
* ```-ntfs``` - consider ntfs alternate file streams
* ```-exec``` - execute for each duplicate {}

### TODO

* integrate cobra cli or hand lex/parse cli arguments
* implement size, name, mtime cli options
* implement unary operators - and, not, or
* integrate cheggaaa/pb or hand build progress bar
* integrate fatih/color or hand colorize output

### Dev Notes

Uses golang concurrency pipeline pattern.

```
generate file list
 -> gather sizes and modes 
  -> hash duplicates concurrently
   -> distill duplicates 
    -> print
```

### Testing

##### Operating Systems

* osx
* linux, centos, debian, ubuntu
* windows

##### File Systems

* apfs
* ext3, ext4
* fat, fat32
* ntfs
* btrfs, cow
* fuse

##### References

* https://golang.org
* https://blog.golang.org/pipelines
* https://github.com/adrianlopezroche/fdupes
* https://github.com/jbruchon/jdupes
* https://github.com/spf13/pflag
* https://github.com/cheggaaa/pb
* https://github.com/vbauerster/mpb
* https://github.com/fatih/color
