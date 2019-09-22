package main

import "stindex/jobs"

func main() {
	// create a Job based on the given arguments
	job := jobs.CreateJob()
	// process input into interim data structure (a Binary Search Tree)
	job.Run()
	// unzip interim results into the final log file
	job.PersistResults()
}
