#!/bin/sh

go run main.go

time ./gfdupes -recurse .
time fdupes -r .

./gfdupes . | sed '/^$/d' | sort
fdupes . | sed '/^$/d' | sort

./gfdupes -recurse . | sed '/^$/d' | sort
fdupes -r . | sed '/^$/d' | sort
