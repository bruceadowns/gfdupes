#!/bin/sh

gfdupes [dir1]
gfdupes [dir1] [dir2]

# show version
gfdupes -version

# show help
gfdupes -h
gfdupes -help

# debug logging to stderr
gfdupes -verbose .

# use md5 hash
gfdupes -hash md5 .

# use sha1 hash
gfdupes -hash sha1 .

# recursively walk the directory for files
gfdupes -recurse .

# consider permissions in file difference
gfdupes -perm .

# do not consider hidden files
gfdupes -nohidden .

# do not consider empty files
gfdupes -noempty .

# do not cross device boundaries
gfdupes -xdev .

# use channel buffer size of 64
gfdupes -buffer 64 .

# follow symlinks
gfdupes -symlinks .

# consider hardlinks
gfdupes -hardlinks .

# prompt for deletion of duplications
gfdupes -delete .

# consider the gunzipped file stats in difference
gfdupes -gzip .

# traverse ntfs alternate file streams
gfdupes -ntfs .

# consider files of size gte to 10mb and lte 20mb
gfdupes -size +10 -size -20m .

# consider files of name eq to foo.txt
gfdupes -name foobar.txt

# consider files of name like foo*.txt
gfdupes -name "foo*.txt"

# consider files of name like *.go and not like *_test.go
gfdupes -name "*.go" ! -name "*_test.go"

# consider files of name like *.java and not like target
gfdupes -name *.java -not -name target

# consider files at a maximum depty of 2
gfdupes -maxdepth 2 .

# consider files modified within 2 days
gfdupes -mtime 2d .

# consider files modified within 1 year
gfdupes -mtime 1y .

# execute rm for each duplication found
gfdupes -exec rm {}
