# stindex¸
A program to sort the strings of various files in a memory-efficient way.

# About
stindex is a small project that I did while learning golang. The challenge was to write a program which would have several files as input. Each of the files would host an unordered list of search terms. The program would ultimatelly create an output file with the ordered list of the search terms of the individual files. An added twist was that the available memory to host the search terms should be variable and could go down to 1. That is, at any time, the program could only have just one term hosted in memory.

# Basic Architectural Tenets
The workhorse data structure is a special Binary Search Tree (BST) and its code lives in package "trees".
The nodes of the tree -and the data atoms in our case- are laid out as the SearchTerm struct
and the rationale is that each separate SearchTerm will also keep track of its own statistics.
BST allows us to lexicographically order search terms as we see them while keeping logical RAM
usage to a minimum. Each insertion needs at most 2 elements in RAM - the one being inserted and the
current tree node that we are comparing with. That way we can work with minimum RAM size of 2. 
Insertions will be usually cheap (where the tree is roughly balanced) - of logn complexity. 
Getting from that tree to the ordered list can be done quickly with an inorder traversal.
Unzipping to the final file can take significantly longer as it is directly proportional to the
number of search terms - O(n) thus I broke the whole process in two separate steps to make the
distinction clear.

Larger logical RAM size speeds us up by allowing reading in more elements from input log files at each
disk access thus minimizing the expensive disk reads. Disk reads is managed by package "readers"
which essentially provides iterator-like access to input files (utilizing Go's channels).

Each separate run of the program is viewed as a job (much akin a batch or a cron job). Job code
lives in package "jobs" and it is laid up so that it taking care of its own configuration, validity
and execution.


## Typical execution (Large input files)
```
stindex$ cat tests_resources/* | wc -l
 16243100
stindex$ ./stindex -outfile=final.log tests_resources/*
2019/07/09 19:34:52 [INFO] RAM Size: 1000000
2019/07/09 19:34:52 [INFO] Input files: [tests_resources/1.log tests_resources/2.log tests_resources/3.log tests_resources/4.log]
2019/07/09 19:34:52 [INFO] Starting processing.
.................
2019/07/09 19:35:11 [INFO] Processed 16243100 search terms. Elapsed time: 18.275413496s
2019/07/09 19:36:15 [INFO] Unzipped to destination file: final.log Elapsed time: 1m4.521703887s
```

## Typical execution (small files)
```
2019/07/09 20:01:34 [INFO] RAM Size: 4
2019/07/09 20:01:34 [INFO] Input files: [tests_resources/1.log tests_resources/2.log tests_resources/3.log]
2019/07/09 20:01:34 [INFO] Starting processing.
......
2019/07/09 20:01:34 [INFO] Processed 20 search terms. Elapsed time: 951.982µs
2019/07/09 20:01:34 [INFO] Unzipped to destination file: out Elapsed time: 439.608µs
stindex$ cat out 
about
aeon
aeon
aeon
beer
beer
beer
beer
bread
curry
deer
fortytwo
locus
lol
one
ship
ship
talk
three
two
```

# How to build
Unzip code into $GOPATH/src
IMPORTANT: Make sure you do not have an stindex folder in $GOPATH/src already.
If you do please move that temporarily somewhere else in order to test this code.
I would host this code under my github account thus removing any chance of colliding
but instructions were not to host this code in public github.
> unzip stindex.zip -d $GOPATH/src
> cd into $GOPATH/src/stindex
> go build

# How to run tests:
> go test ./...

# How to get help/see usage:
> ./stindex -h
