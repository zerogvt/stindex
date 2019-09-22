package readers

import (
	"bufio"
	"os"
	"strings"
)

//linesGenerator is a Python generator-like file reader.
//Input:  A file path
//Output: An unbuffered channel emitting a file's line at each read
func linesGenerator(filename string) chan string {
	c := make(chan string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			close(c)
			return
		}
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(c)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			c <- line
		}
	}()
	return c
}

// ReadFilesSerially reads a set of files line by line serially file after file.
// Thus given an array of files {file1, file2, file3} it
// will read all lines off file1, then all lines off file2 and so on
// till it exhausts all files. It returns an unbuffered channel that will
// return a line at a time.
func ReadFilesSerially(files []string) chan string {
	// return channel
	sink := make(chan string)
	// create a channel for each file
	srcs := make([]chan string, len(files))
	for i, f := range files {
		srcs[i] = linesGenerator(f)
	}
	go func() {
		for _, src := range srcs {
			// read all lines off current src file
			for {
				line, ok := <-src
				if ok {
					sink <- line
				} else {
					// when current src is exhausted move on
					break
				}
			}
		}
		close(sink)
		return
	}()
	return sink
}
