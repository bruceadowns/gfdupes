#!/bin/sh

$ ./fdupes -h
Usage: fdupes [options] DIRECTORY...

 -r --recurse     	for every directory given follow subdirectories
                  	encountered within
 -R --recurse:    	for each directory given after this option follow
                  	subdirectories encountered within (note the ':' at
                  	the end of the option, manpage for more details)
 -s --symlinks    	follow symlinks
 -H --hardlinks   	normally, when two or more files point to the same
                  	disk area they are treated as non-duplicates; this
                  	option will change this behavior
 -n --noempty     	exclude zero-length files from consideration
 -A --nohidden    	exclude hidden files from consideration
 -f --omitfirst   	omit the first file in each set of matches
 -1 --sameline    	list each set of matches on a single line
 -S --size        	show size of duplicate files
 -m --summarize   	summarize dupe information
 -q --quiet       	hide progress indicator
 -d --delete      	prompt user for files to preserve and delete all
                  	others; important: under particular circumstances,
                  	data may be lost when using this option together
                  	with -s or --symlinks, or when specifying a
                  	particular directory more than once; refer to the
                  	fdupes documentation for additional information
 -N --noprompt    	together with --delete, preserve the first file in
                  	each set of duplicates and delete the rest without
                  	prompting the user
 -I --immediate   	delete duplicates as they are encountered, without
                  	grouping into sets; implies --noprompt
 -p --permissions 	don't consider files with different owner/group or
                  	permission bits as duplicates
 -o --order=BY    	select sort order for output and deleting; by file
                  	modification time (BY='time'; default), status
                  	change time (BY='ctime'), or filename (BY='name')
 -i --reverse     	reverse order while sorting
 -v --version     	display fdupes version
 -h --help        	display this help message
