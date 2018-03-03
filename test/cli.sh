#!/bin/sh

go run main.go

gfdupes [dir1]
gfdupes [dir1] [dir2]

gfdupes -recurse [dir1]
gfdupes -perm [dir1]

gfdupes -hash md5 [dir1]
gfdupes -hash sha1 [dir1]

gfdupes -help
