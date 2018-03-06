#!/bin/sh

go run main.go

time ./gfdupes -recurse .
time ~/dev/scratch/fdupes/fdupes -r .

./gfdupes . | sed '/^$/d' | sort
~/dev/scratch/fdupes/fdupes . | sed '/^$/d' | sort

./gfdupes -recurse . | sed '/^$/d' | sort
~/dev/scratch/fdupes/fdupes -r . | sed '/^$/d' | sort
